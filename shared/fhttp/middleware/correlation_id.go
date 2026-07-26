package middleware

import (
	"context"
	"net/http"

	"github.com/cchristian77/wallet-service/util/constant"
	"github.com/cchristian77/wallet-service/util/logger"
	"github.com/google/uuid"
)

// CorrelationID - Middleware to add centralized requestID per incoming to context
// if "Correlation-ID" is empty then add this in header and context
// if "Idempotency-Key" exists then add this context
func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := r.Header.Get(constant.XCorrelationIDKey)
		if requestID == "" {
			logger.Debug(ctx, "Creating new request id.")
			requestID = uuid.New().String()

			newRequest := r.WithContext(ctx)
			newRequest.Header.Set(constant.XCorrelationIDKey, requestID)
		}

		ctx = context.WithValue(ctx, constant.XCorrelationIDKey, requestID)
		logger.Debug(ctx, "CorrelationID: %v", requestID)

		idempotencyKey := r.Header.Get(constant.XIdempotencyKey)
		if idempotencyKey != "" {
			logger.Debug(ctx, "Idempotency-Key found from the header : %s", idempotencyKey)
			ctx = context.WithValue(ctx, constant.XIdempotencyKey, idempotencyKey)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
