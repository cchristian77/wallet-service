package repository

import (
	"context"

	"github.com/cchristian77/wallet-service/domain"
)

func (r *repo) FindTransactionLedgersByTransactionID(ctx context.Context, transactionID uint64) ([]*domain.TransactionLedger, error) {
	panic("implement me")
}

func (r *repo) CreateTransactionLedger(ctx context.Context, data *domain.TransactionLedger) (*domain.TransactionLedger, error) {
	panic("implement me")
}
