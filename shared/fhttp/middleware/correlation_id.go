package middleware

import (
	"context"
	"net/http"

	"github.com/cchristian77/wallet-service/util/constant"
	"github.com/cchristian77/wallet-service/util/logger"
	"github.com/google/uuid"
)

// CorrelationID attaches Correlation-ID (generated if missing) and Idempotency-Key
// (when present) to the request context and response headers.
func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := r.Header.Get(constant.XCorrelationIDKey)
		if requestID == "" {
			logger.Debug(ctx, "Creating new request id.")
			requestID = uuid.New().String()
		}

		ctx = context.WithValue(ctx, constant.XCorrelationIDKey, requestID)
		logger.Debug(ctx, "CorrelationID: %v", requestID)

		idempotencyKey := r.Header.Get(constant.XIdempotencyKey)
		if idempotencyKey != "" {
			ctx = context.WithValue(ctx, constant.XIdempotencyKey, idempotencyKey)
			w.Header().Set(constant.XIdempotencyKey, idempotencyKey)
			logger.Debug(ctx, "Idempotency-Key from the header : %s", idempotencyKey)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
