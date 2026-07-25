package repository

import (
	"context"

	"github.com/cchristian77/wallet-service/domain"
)

func (r *repo) FindWalletByID(ctx context.Context, id uint64) (*domain.Wallet, error) {
	panic("implement me")
}

func (r *repo) FindWalletByUserID(ctx context.Context, userID uint64) (*domain.Wallet, error) {
	panic("implement me")
}

func (r *repo) FindWalletByIDForUpdate(ctx context.Context, id uint64) (*domain.Wallet, error) {
	panic("implement me")
}

func (r *repo) UpdateWalletBalance(ctx context.Context, walletID uint64, balance int64) error {
	panic("implement me")
}
