package memstore

import (
	"context"
	"time"

	"github.com/cchristian77/wallet-service/domain"
	sharedErrs "github.com/cchristian77/wallet-service/shared/errors"
	"github.com/cchristian77/wallet-service/util/logger"
)

func (s *MemStore) FindTransactionLedgersByTransactionID(ctx context.Context, transactionID uint64) ([]*domain.TransactionLedger, error) {
	s.rLock(ctx)
	defer s.rUnlock(ctx)

	ids := s.ledgersByTxn[transactionID]
	result := make([]*domain.TransactionLedger, 0, len(ids))
	for _, id := range ids {
		if l, ok := s.ledgers[id]; ok {
			result = append(result, cloneLedger(l))
		}
	}

	return result, nil
}

func (s *MemStore) CreateTransactionLedger(ctx context.Context, data *domain.TransactionLedger) (*domain.TransactionLedger, error) {
	s.wLock(ctx)
	defer s.wUnlock(ctx)

	if _, ok := s.transactions[data.TransactionID]; !ok {
		logger.Error(ctx, "[MEMSTORE] Failed on create transaction ledger: transaction not found")
		return nil, sharedErrs.NotFoundErr
	}
	if _, ok := s.wallets[data.WalletID]; !ok {
		logger.Error(ctx, "[MEMSTORE] Failed on create transaction ledger: wallet not found")
		return nil, sharedErrs.NotFoundErr
	}

	now := time.Now().UTC()
	created := cloneLedger(data)
	created.ID = s.nextLedgerID
	s.nextLedgerID++
	if created.CreatedAt.IsZero() {
		created.CreatedAt = now
	}
	created.UpdatedAt = now

	s.ledgers[created.ID] = created
	s.ledgersByTxn[created.TransactionID] = append(s.ledgersByTxn[created.TransactionID], created.ID)

	return cloneLedger(created), nil
}
