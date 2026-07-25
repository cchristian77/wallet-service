package repository

import (
	"context"

	"github.com/cchristian77/wallet-service/domain"
	sharedErrs "github.com/cchristian77/wallet-service/shared/errors"
	"github.com/cchristian77/wallet-service/shared/external/database"
	"github.com/cchristian77/wallet-service/util/logger"

	"gorm.io/gorm/clause"
)

func (r *repo) FindWalletByID(ctx context.Context, id uint64) (*domain.Wallet, error) {
	var result *domain.Wallet

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		First(&result, id).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find wallet by id : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) FindWalletByUserID(ctx context.Context, userID uint64) (*domain.Wallet, error) {
	var result *domain.Wallet

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&result).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find wallet by user id : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) FindWalletByIDForUpdate(ctx context.Context, id uint64) (*domain.Wallet, error) {
	var result *domain.Wallet

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&result, id).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on find wallet by id for update : %v", err)

		return nil, sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return result, nil
}

func (r *repo) UpdateWalletBalance(ctx context.Context, id uint64, balance int64) error {
	var result *domain.Wallet

	db, _ := database.ConnFromCtx(ctx, r.DB)

	err := db.WithContext(ctx).
		Model(&result).
		Where("id = ?", id).
		Update("balance", balance).
		Error
	if err != nil {
		logger.Error(ctx, "[REPOSITORY] Failed on update wallet balance: %v", err)

		return sharedErrs.NewRepositoryErr(err, "%s", err.Error())
	}

	return nil
}
