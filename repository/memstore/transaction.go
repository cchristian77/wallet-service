package memstore

import (
	"context"
	"time"

	"github.com/cchristian77/wallet-service/domain"
	"github.com/cchristian77/wallet-service/domain/enums"
	sharedErrs "github.com/cchristian77/wallet-service/shared/errors"
	"github.com/cchristian77/wallet-service/util/logger"
)

func (s *MemStore) FindTransactionByTransactionID(ctx context.Context, transactionID string) (*domain.Transaction, error) {
	s.rLock(ctx)
	defer s.rUnlock(ctx)

	id, ok := s.txByKey[transactionID]
	if !ok {
		logger.Error(ctx, "[MEMSTORE] Failed on find transaction by transaction id : not found")
		return nil, sharedErrs.NotFoundErr
	}

	t, ok := s.transactions[id]
	if !ok {
		logger.Error(ctx, "[MEMSTORE] Failed on find transaction by transaction id : not found")
		return nil, sharedErrs.NotFoundErr
	}

	return cloneTransaction(t), nil
}

func (s *MemStore) CreateTransaction(ctx context.Context, data *domain.Transaction) (*domain.Transaction, error) {
	s.wLock(ctx)
	defer s.wUnlock(ctx)

	if _, exists := s.txByKey[data.TransactionID]; exists {
		logger.Error(ctx, "[MEMSTORE] Failed on create transaction: conflict")
		return nil, sharedErrs.ConflictErr
	}

	now := time.Now().UTC()
	created := cloneTransaction(data)
	created.ID = s.nextTxnID
	s.nextTxnID++
	if created.CreatedAt.IsZero() {
		created.CreatedAt = now
	}
	created.UpdatedAt = now

	s.transactions[created.ID] = created
	s.txByKey[created.TransactionID] = created.ID

	return cloneTransaction(created), nil
}

func (s *MemStore) UpdateTransactionStatus(ctx context.Context, id uint64, status string) error {
	s.wLock(ctx)
	defer s.wUnlock(ctx)

	t, ok := s.transactions[id]
	if !ok {
		logger.Error(ctx, "[MEMSTORE] Failed on update transaction status: not found")
		return sharedErrs.NotFoundErr
	}

	t.Status = enums.TransactionStatus(status)
	t.UpdatedAt = time.Now().UTC()

	return nil
}
