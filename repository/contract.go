package repository

import (
	"context"

	"github.com/cchristian77/wallet-service/domain"
)

//go:generate mockgen -package repository -source=contract.go -destination=mock_repository.go *

// Repository defines data access methods for disbursement domains.
type Repository interface {
	WalletRepository
	TransactionRepository
	TransactionLedgerRepository
}

type WalletRepository interface {
	FindWalletByID(ctx context.Context, id uint64) (*domain.Wallet, error)
	FindWalletByUserID(ctx context.Context, userID uint64) (*domain.Wallet, error)
	FindWalletByIDForUpdate(ctx context.Context, id uint64) (*domain.Wallet, error)
	UpdateWalletBalance(ctx context.Context, walletID uint64, balance int64) error
}

type TransactionRepository interface {
	FindTransactionByTransactionID(ctx context.Context, transactionID string) (*domain.Transaction, error)
	CreateTransaction(ctx context.Context, data *domain.Transaction) (*domain.Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id uint64, status string) error
}

type TransactionLedgerRepository interface {
	FindTransactionLedgersByTransactionID(ctx context.Context, transactionID uint64) ([]*domain.TransactionLedger, error)
	CreateTransactionLedger(ctx context.Context, data *domain.TransactionLedger) (*domain.TransactionLedger, error)
}
