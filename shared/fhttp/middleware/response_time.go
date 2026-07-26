package middleware

import (
	"net/http"
	"time"

	"github.com/cchristian77/wallet-service/util/logger"
)

const timeFormat = "2006-01-02 15:04:05.000000"

// ResponseTime - Middleware to log the response time for each request made
func ResponseTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		startTime := time.Now()

		logger.Info(ctx, "Request made at: %v url %v", startTime.Format(timeFormat), r.URL)
		next.ServeHTTP(w, r)
		logger.Info(ctx, "Request duration: %v url %v", time.Since(startTime).String(), r.URL)
	})
}
