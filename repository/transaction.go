package repository

import (
	"context"

	"github.com/cchristian77/wallet-service/domain"
)

func (r *repo) FindTransactionByTransactionID(ctx context.Context, transactionID string) (*domain.Transaction, error) {
	panic("implement me")
}

func (r *repo) CreateTransaction(ctx context.Context, data *domain.Transaction) (*domain.Transaction, error) {
	panic("implement me")
}

func (r *repo) UpdateTransactionStatus(ctx context.Context, id uint64, status string) error {
	panic("implement me")
}
