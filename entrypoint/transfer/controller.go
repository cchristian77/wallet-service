package transfer

import (
	"encoding/json"
	"net/http"

	"github.com/cchristian77/wallet-service/request"
	transactionLedger "github.com/cchristian77/wallet-service/service/transaction_ledger"
	sharedErrs "github.com/cchristian77/wallet-service/shared/errors"
	"github.com/cchristian77/wallet-service/shared/fhttp"
	"github.com/cchristian77/wallet-service/util"
	"github.com/cchristian77/wallet-service/util/constant"
)

// Controller manages the transfer/disbursement operations.
type Controller struct {
	transactionLedger transactionLedger.Service
}

func (c *Controller) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /transfers/v1", fhttp.AppHandler(c.Transfer))
}

func (c *Controller) Transfer(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	idempotencyKey := constant.IdempotencyKeyFromCtx(ctx)
	if idempotencyKey == "" {
		idempotencyKey = r.Header.Get(constant.XIdempotencyKey)
	}

	if idempotencyKey == "" {
		return nil, fhttp.NewErrorResponse(
			http.StatusBadRequest,
			sharedErrs.ErrKindValidation.String(),
			"Idempotency-Key header is required",
		)
	}

	var input request.Transfer
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusUnprocessableEntity,
			sharedErrs.ErrKindInvalidRequest.String(),
			"Invalid request body",
		)
	}

	input.IdempotencyKey = idempotencyKey
	if err := util.Validate(input); err != nil {
		return nil, err
	}

	data, err := c.transactionLedger.Transfer(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{
		Status:  http.StatusOK,
		Message: "OK",
		Data:    data,
	}, nil
}
