package database

import (
	"context"

	"gorm.io/gorm"
)

type keys string

const txKey keys = "DBTRX"

func ConnFromContext(ctx context.Context, defaults ...*gorm.DB) (*gorm.DB, bool) {
	tx, ok := GetTxFromContext(ctx)
	if ok {
		return tx, ok
	}

	if len(defaults) == 0 {
		return nil, false
	}

	return defaults[0], true
}

func InitTx(ctx context.Context, db *gorm.DB) (context.Context, *gorm.DB) {
	dbTx := db.WithContext(ctx).Begin()
	return context.WithValue(ctx, txKey, dbTx), dbTx
}

func GetTxFromContext(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txKey).(*gorm.DB)
	return tx, ok
}
