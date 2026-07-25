package transaction_ledger

import (
	"context"

	"github.com/cchristian77/wallet-service/request"
	"github.com/cchristian77/wallet-service/response"
)

func (b *base) Transfer(ctx context.Context, input *request.Transfer) (*response.Transfer, error) {
	panic("implement me")
}
