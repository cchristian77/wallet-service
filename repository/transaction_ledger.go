package repository

import (
	"context"

	"github.com/cchristian77/wallet-service/domain"
	sharedErrs "github.com/cchristian77/wallet-service/shared/errors"
	"github.com/cchristian77/wallet-service/shared/external/database"
	"github.com/cchristian77/wallet-service/util/logger"

	"gorm.io/gorm/clause"
)

func (r *repo) FindTransactionLedgersByTransactionID(ctx context.Context, transactionID uint64) ([]*domain.TransactionLedger, error) {
	var result []*domain.TransactionLedger

	db, _ := database.ConnFromContext(ctx, r.DB)

	err := db.WithContext(ctx).
		Where("transaction_id = ?", transactionID).
		Find(&result).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find transaction ledgers by transaction id : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) CreateTransactionLedger(ctx context.Context, data *domain.TransactionLedger) (*domain.TransactionLedger, error) {
	db, _ := database.ConnFromContext(ctx, r.DB)

	err := db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(&data).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on create transaction ledger: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return data, nil
}
