package transfer

import (
	"context"
	"fmt"

	"github.com/cchristian77/wallet-service/repository"
	transactionLedger "github.com/cchristian77/wallet-service/service/transaction_ledger"
	"github.com/cchristian77/wallet-service/util/logger"
	"gorm.io/gorm"
)

// NewController initializes a new Controller instance.
func NewController(ctx context.Context, repository repository.Repository, writerDB *gorm.DB) (*Controller, error) {
	transactionLedgerService, err := transactionLedger.NewService(repository, writerDB)
	if err != nil {
		logger.L().Fatal(fmt.Sprintf("transaction_ledger service initialization error: %v", err))
	}

	return &Controller{transactionLedger: transactionLedgerService}, nil
}
