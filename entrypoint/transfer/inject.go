package transfer

import (
	"context"
	"fmt"

	"github.com/cchristian77/wallet-service/repository"
	transactionLedger "github.com/cchristian77/wallet-service/service/transaction_ledger"
	"github.com/cchristian77/wallet-service/shared/external/database"
	"github.com/cchristian77/wallet-service/util/logger"
)

// NewController initializes a new Controller instance.
func NewController(ctx context.Context, repository repository.Repository, transactor database.Transactor) (*Controller, error) {
	transactionLedgerService, err := transactionLedger.NewService(repository, transactor)
	if err != nil {
		logger.L().Fatal(fmt.Sprintf("transaction_ledger service initialization error: %v", err))
	}

	return &Controller{transactionLedger: transactionLedgerService}, nil
}
