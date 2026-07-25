package constant

import "context"

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	correlationIDKey = contextKey("Correlation-ID")
	idempotencyKey   = contextKey("Idempotency-Key")
)

var (
	XCorrelationIDKey = correlationIDKey.String()
	XIdempotencyKey   = idempotencyKey.String()
)

func CorrelationIDFromCtx(ctx context.Context) string {
	correlationID, ok := ctx.Value(XCorrelationIDKey).(string)
	if !ok {
		return ""
	}

	return correlationID
}

func IdempotencyKeyFromCtx(ctx context.Context) string {
	key, ok := ctx.Value(XIdempotencyKey).(string)
	if !ok {
		return ""
	}

	return key
}
