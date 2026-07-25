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

var (
	transactionAlreadyProcessedErr = sharedErrs.New(sharedErrs.ErrKindConflict, "Transaction already processed")
	insufficientBalanceErr         = sharedErrs.NewBusinessValidationErr("Insufficient balance")
	sameWalletErr                  = sharedErrs.NewBusinessValidationErr("Cannot transfer to the same wallet")
)

func (b *base) Transfer(ctx context.Context, input *request.Transfer) (*response.Transfer, error) {
	logger.Info(ctx, "Transfer with req: %v", input)

	if input.From == input.To {
		logger.Warn(ctx, "transfer rejected: same wallet id %d", input.From)
		return nil, sameWalletErr
	}

	var err error

	tCtx, tx := database.InitTx(ctx, b.writeDB)
	defer func() {
		if err = tx.Rollback().Error; err != nil && !errors.Is(err, sql.ErrTxDone) {
			logger.Error(ctx, "Repository Error on executing b.Transfer: ROLLBACK TXN: %v", err)
		}
	}()

	logger.Debug(ctx, "[TRANSFER] checking whether transaction with idempotency key %s already exists ...", input.IdempotencyKey)
	transactionExists, err := b.repository.FindTransactionByTransactionID(tCtx, input.IdempotencyKey)
	if err != nil && !errors.Is(err, sharedErrs.NotFoundErr) {
		return nil, err
	}
	if transactionExists != nil {
		logger.Warn(ctx, "[TRANSFER] transaction %s with idempotency key %s already processed",
			transactionExists.ID, transactionExists.TransactionID)
		return nil, transactionAlreadyProcessedErr
	}

	logger.Info(ctx, "[TRANSFER] creating transaction with idempotency key %s ...", input.IdempotencyKey)
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

	logger.Info(ctx, "[TRANSFER] locking wallets %d and %d for update ...", firstID, secondID)
	firstWallet, err := b.repository.FindWalletByIDForUpdate(tCtx, firstID)
	if err != nil {
		if errors.Is(err, sharedErrs.NotFoundErr) {
			return nil, sharedErrs.New(sharedErrs.ErrKindDataNotFound, "Wallet id %d not found", firstID)
		}
		return nil, err
	}
	secondWallet, err := b.repository.FindWalletByIDForUpdate(tCtx, secondID)
	if err != nil {
		if errors.Is(err, sharedErrs.NotFoundErr) {
			return nil, sharedErrs.New(sharedErrs.ErrKindDataNotFound, "Wallet id %d not found", secondID)
		}
		return nil, err
	}

	var fromWallet, toWallet *domain.Wallet
	if firstWallet.ID == input.From {
		fromWallet, toWallet = firstWallet, secondWallet
	} else {
		fromWallet, toWallet = secondWallet, firstWallet
	}

	logger.Info(ctx, "[TRANSFER] source wallet %d initial balance: %d to wallet %d initial balance: %d with transfer amount: %d",
		fromWallet.ID, fromWallet.Balance, toWallet.ID, toWallet.Balance, input.Amount)
	if fromWallet.Balance < input.Amount {
		logger.Warn(ctx, "[TRANSFER] wallet %d has insufficient balance %d to from transfer amount %d",
			fromWallet.ID, fromWallet.Balance, input.Amount)
		return nil, insufficientBalanceErr
	}

	logger.Info(ctx, "[TRANSFER] creating CREDIT ledger for wallet %d with amount %d ...", toWallet.ID, input.Amount)
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

	logger.Info(ctx, "[TRANSFER] creating DEBIT ledger for wallet %d amount %d ...", fromWallet.ID, input.Amount)
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
	toWallet.Balance += input.Amount

	logger.Info(ctx, "[TRANSFER]  source wallet %s to final balance %d...", fromWallet.ID, fromWallet.Balance)
	if err = b.repository.UpdateWalletBalance(tCtx, fromWallet.ID, fromWallet.Balance); err != nil {
		return nil, err
	}

	logger.Info(ctx, "[TRANSFER] updating designated wallet %s to final balance %d...", fromWallet.ID, fromWallet.Balance)
	if err = b.repository.UpdateWalletBalance(tCtx, toWallet.ID, toWallet.Balance); err != nil {
		return nil, err
	}

	logger.Info(ctx, "[TRANSFER] updating transaction %d to SUCCESS ...", transaction.ID)
	if err = b.repository.UpdateTransactionStatus(tCtx, transaction.ID, enums.TransactionStatusSuccess.String()); err != nil {
		return nil, err
	}
	transaction.Status = enums.TransactionStatusSuccess

	if err = tx.Commit().Error; err != nil {
		logger.Error(ctx, "Repository Error on executing b.Transfer: COMMIT TXN: %v", err)
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
