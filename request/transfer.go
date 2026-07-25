package request

type Transfer struct {
	From           uint64 `json:"from" validate:"required"`
	To             uint64 `json:"to" validate:"required"`
	Amount         int64  `json:"amount" validate:"required,gt=0"`
	IdempotencyKey string `json:"-" validate:"required,max=255"` // from Idempotency-Key header
}
