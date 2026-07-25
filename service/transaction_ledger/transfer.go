package transaction_ledger

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cchristian77/wallet-service/domain"
	"github.com/cchristian77/wallet-service/domain/enums"
	"github.com/cchristian77/wallet-service/request"
	"github.com/cchristian77/wallet-service/response"
	sharedErrs "github.com/cchristian77/wallet-service/shared/errors"
	"github.com/cchristian77/wallet-service/shared/external/database"
	"github.com/cchristian77/wallet-service/util/logger"
)

func (b *base) Transfer(ctx context.Context, input *request.Transfer) (*response.Transfer, error) {
	logger.Info(ctx, "Transfer with req: %v", input)

	if input.From == input.To {
		return nil, SameWalletErr
	}

	var err error

	tCtx, tx := database.InitTx(ctx, b.writeDB)
	defer func() {
		if err = tx.Rollback().Error; err != nil && !errors.Is(err, sql.ErrTxDone) {
			logger.Error(ctx, "Repository Error on executing b.Claim: ROLLBACK TXN: %v", err)
		}
	}()

	transactionExists, err := b.repository.FindTransactionByTransactionID(tCtx, input.IdempotencyKey)
	if err != nil && !errors.Is(err, sharedErrs.NotFoundErr) {
		return nil, err
	}
	if transactionExists != nil {
		return nil, TransactionAlreadyProcessedErr
	}

	transaction, err := b.repository.CreateTransaction(tCtx, &domain.Transaction{
		TransactionID: input.IdempotencyKey,
		Status:        enums.TransactionStatusPending,
	})
	if err != nil {
		return nil, err
	}

	// Lock wallets in ascending ID order to avoid deadlocks (A→B vs B→A).
	firstID, secondID := input.From, input.To
	if firstID > secondID {
		firstID, secondID = secondID, firstID
	}

	firstWallet, err := b.repository.FindWalletByIDForUpdate(tCtx, firstID)
	if err != nil {
		if errors.Is(err, sharedErrs.NotFoundErr) {
			return nil, sharedErrs.New(sharedErrs.ErrKindDataNotFound, "Wallet with id %d not found", firstID)
		}
		return nil, err
	}
	secondWallet, err := b.repository.FindWalletByIDForUpdate(tCtx, secondID)
	if err != nil {
		if errors.Is(err, sharedErrs.NotFoundErr) {
			return nil, sharedErrs.New(sharedErrs.ErrKindDataNotFound, "Wallet with id %d not found", firstID)
		}
		return nil, err
	}

	var fromWallet, toWallet *domain.Wallet
	if firstWallet.ID == input.From {
		fromWallet, toWallet = firstWallet, secondWallet
	} else {
		fromWallet, toWallet = secondWallet, firstWallet
	}

	if fromWallet.Balance < input.Amount {
		return nil, InsufficientBalanceErr
	}

	// CREDIT recipient
	_, err = b.repository.CreateTransactionLedger(tCtx, &domain.TransactionLedger{
		TransactionID: transaction.ID,
		WalletID:      toWallet.ID,
		Ledger:        enums.TransactionLedgerCredit,
		Reference:     enums.TransactionLedgerRefTransfer,
		Amount:        input.Amount,
	})
	if err != nil {
		return nil, err
	}

	// DEBIT sender
	_, err = b.repository.CreateTransactionLedger(tCtx, &domain.TransactionLedger{
		TransactionID: transaction.ID,
		WalletID:      fromWallet.ID,
		Ledger:        enums.TransactionLedgerDebit,
		Reference:     enums.TransactionLedgerRefTransfer,
		Amount:        input.Amount,
	})
	if err != nil {
		return nil, err
	}

	fromWallet.Balance -= input.Amount
	if err = b.repository.UpdateWalletBalance(tCtx, fromWallet.ID, fromWallet.Balance); err != nil {
		return nil, err
	}

	toWallet.Balance += input.Amount
	if err = b.repository.UpdateWalletBalance(tCtx, toWallet.ID, toWallet.Balance); err != nil {
		return nil, err
	}

	if err = b.repository.UpdateTransactionStatus(tCtx, transaction.ID, enums.TransactionStatusSuccess.String()); err != nil {
		return nil, err
	}
	transaction.Status = enums.TransactionStatusSuccess

	if err = tx.Commit().Error; err != nil {
		logger.Error(ctx, "Transfer: COMMIT TXN: %v", err)
		return nil, err
	}

	return &response.Transfer{
		TransactionID: transaction.TransactionID,
		FromBalance:   fromWallet.Balance,
		ToBalance:     toWallet.Balance,
		Amount:        input.Amount,
		Status:        transaction.Status.String(),
	}, nil
}
