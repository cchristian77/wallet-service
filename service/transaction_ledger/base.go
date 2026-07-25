package transaction_ledger

import (
	"context"

	"github.com/cchristian77/wallet-service/repository"
	"github.com/cchristian77/wallet-service/request"
	"github.com/cchristian77/wallet-service/response"
	sharedErrs "github.com/cchristian77/wallet-service/shared/errors"
	"gorm.io/gorm"
)

type Service interface {
	Transfer(ctx context.Context, input *request.Transfer) (*response.Transfer, error)
}

type base struct {
	repository repository.Repository
	writeDB    *gorm.DB
}

func NewService(repository repository.Repository, writeDB *gorm.DB) (Service, error) {
	return &base{
		repository: repository,
		writeDB:    writeDB,
	}, nil
}

var (
	TransactionAlreadyProcessedErr = sharedErrs.New(sharedErrs.ErrKindConflict, "Transaction already processed")
	InsufficientBalanceErr         = sharedErrs.NewBusinessValidationErr("Insufficient balance")
	SameWalletErr                  = sharedErrs.NewBusinessValidationErr("Cannot transfer to the same wallet")
)
