package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/cchristian77/wallet-service/util/config"
	"github.com/cchristian77/wallet-service/util/logger"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		healthcheckPath := "/healthcheck"
		if r.URL.Path == healthcheckPath {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()

		if r.Body == nil {
			logger.Info(ctx, "Received Request from URL: %v, with empty body", r.URL.Path)
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil && err != io.EOF {
			logger.Error(ctx, "Error reading request body: %v", err)
		} else {
			logger.Info(ctx, "Received Request from URL: %v, with body: %v", r.URL.Path, string(body))
		}
		_ = r.Body.Close()

		r.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	})
}

func LogResponse(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if config.Env() != nil && config.Env().App.Env == "production" {
			h.ServeHTTP(w, r)
			return
		}

		rec := &loggingResponseWriter{ResponseWriter: w, Status: http.StatusOK}
		h.ServeHTTP(rec, r)

		statusCode := rec.Status
		logMsg := fmt.Sprintf("Received HTTP Status Code: %d, with Response: %s", statusCode, string(rec.Body))

		switch {
		case statusCode >= 500:
			logger.Error(ctx, "%s", logMsg)
		case statusCode >= 400:
			logger.Warn(ctx, "%s", logMsg)
		case statusCode >= 300:
			logger.Info(ctx, "%s", logMsg)
		default:
			logger.Debug(ctx, "%s", logMsg)
		}
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	Status int
	Body   []byte
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.Status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	res, err := lrw.ResponseWriter.Write(b)
	if err == nil {
		lrw.Body = append(lrw.Body, b...)
	}
	return res, err
}
