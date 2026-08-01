package transaction_ledger

import (
	"context"

	"github.com/cchristian77/wallet-service/repository"
	"github.com/cchristian77/wallet-service/request"
	"github.com/cchristian77/wallet-service/response"
	"github.com/cchristian77/wallet-service/shared/external/database"
)

// Service defines application logic for wallet-to-wallet transfer (disbursement).
type Service interface {
	Transfer(ctx context.Context, input *request.Transfer) (*response.Transfer, error)
}

type base struct {
	repository repository.Repository
	transactor database.Transactor
}

// NewService initializes a new transaction ledger Service instance.
func NewService(repository repository.Repository, transactor database.Transactor) (Service, error) {
	return &base{
		repository: repository,
		transactor: transactor,
	}, nil
}
