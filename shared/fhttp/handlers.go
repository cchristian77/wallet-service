package fhttp

import (
	"context"
	"net/http"

	"github.com/cchristian77/wallet-service/util/constant"
)

// AppHandler wraps controller functions that return (*Response, error).
type AppHandler func(*http.Request) (*Response, error)

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp, err := fn(r)

	if err != nil {
		WriteErrorResponse(ctx, err, w)
		return
	}

	WriteHTTPResponse(ctx, resp, w)
}

func setHeaders(ctx context.Context, headers HTTPHeaders, w http.ResponseWriter) {
	if headers == nil {
		headers = HTTPHeaders{}
	}

	if _, ok := headers[ContentTypeKey]; !ok {
		w.Header().Set(ContentTypeKey, string(ContentTypeJSON))
	}

	if val, ok := ctx.Value(constant.XCorrelationIDKey).(string); ok {
		w.Header().Set(constant.XCorrelationIDKey, val)
	}

	for key, val := range headers {
		w.Header().Set(key, val)
	}
}

func DefaultHealthCheckHandler(_ *http.Request) (*Response, error) {
	return &Response{
		Data:   "Service is Running.",
		Status: http.StatusOK,
	}, nil
}
