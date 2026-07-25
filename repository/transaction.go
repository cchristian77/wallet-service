package repository

import (
	"context"

	"github.com/cchristian77/wallet-service/domain"
	sharedErrs "github.com/cchristian77/wallet-service/shared/errors"
	"github.com/cchristian77/wallet-service/shared/external/database"
	"github.com/cchristian77/wallet-service/util/logger"

	"gorm.io/gorm/clause"
)

func (r *repo) FindTransactionByTransactionID(ctx context.Context, transactionID string) (*domain.Transaction, error) {
	var result *domain.Transaction

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Where("transaction_id = ?", transactionID).
		First(&result).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find transaction by transaction id : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) CreateTransaction(ctx context.Context, data *domain.Transaction) (*domain.Transaction, error) {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(&data).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on create transaction: %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return data, nil
}

func (r *repo) UpdateTransactionStatus(ctx context.Context, id uint64, status string) error {
	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Model(&domain.Transaction{}).
		Where("id = ?", id).
		Update("status", status).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on update transaction status: %v", err)

		return sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return nil
}
