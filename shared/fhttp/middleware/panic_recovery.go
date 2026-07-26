package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/cchristian77/wallet-service/shared/fhttp"
	"github.com/cchristian77/wallet-service/util/logger"
)

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if p := recover(); p != nil {
				ctx := r.Context()
				logger.Error(ctx, "Panic: %v, Panic stacktrace: \n %v", p, string(debug.Stack()))
				fhttp.WriteErrorResponse(ctx, fmt.Errorf("Internal Server Error. Panic Error : %v", p), w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
