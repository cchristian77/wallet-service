package database

import (
	"context"

	"gorm.io/gorm"
)

// Tx is a unit-of-work handle. Commit persists staged changes; Rollback discards them.
// After either succeeds, further Commit/Rollback calls should return sql.ErrTxDone.
type Tx interface {
	Commit() error
	Rollback() error
}

// Transactor starts a transaction and stores it on the returned context so
// repository methods can participate in the same unit of work.
type Transactor interface {
	Begin(ctx context.Context) (context.Context, Tx)
}

// GormTransactor starts GORM/Postgres transactions.
type GormTransactor struct {
	DB *gorm.DB
}

// NewGormTransactor wraps a *gorm.DB as a Transactor.
func NewGormTransactor(db *gorm.DB) *GormTransactor {
	return &GormTransactor{DB: db}
}

func (g *GormTransactor) Begin(ctx context.Context) (context.Context, Tx) {
	ctx, dbTx := InitTx(ctx, g.DB)
	return ctx, &gormTx{db: dbTx}
}

type gormTx struct {
	db *gorm.DB
}

func (t *gormTx) Commit() error {
	return t.db.Commit().Error
}

func (t *gormTx) Rollback() error {
	return t.db.Rollback().Error
}
