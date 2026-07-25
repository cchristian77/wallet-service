package fhttp

import (
	"context"
	"encoding/json"
	"net/http"

	sharedErrs "github.com/cchristian77/wallet-service/shared/errors"
	"github.com/cchristian77/wallet-service/util/constant"
	"github.com/cchristian77/wallet-service/util/logger"
)

type ContentType string

const (
	ContentTypeKey  = "Content-Type"
	ContentTypeJSON = ContentType("application/json")
)

type (
	HTTPHeaders map[string]string
	MetaData    map[string]any
)

func (h HTTPHeaders) Add(key, val string) {
	h[key] = val
}

func (m MetaData) Add(key string, val any) {
	m[key] = val
}

type Response struct {
	Data     any    `json:"data"`
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Headers  HTTPHeaders
	MetaData MetaData
}

type errorResponse struct {
	Status        int            `json:"status"`
	Message       string         `json:"message"`
	Kind          string         `json:"kind"`
	OptionalData  []OptionalData `json:"optional_data"`
	CorrelationID string         `json:"correlation_id"`
}

type OptionalData struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func (e errorResponse) Error() string {
	return e.Message
}

func NewErrorResponse(status int, kind string, message string, optionalData ...OptionalData) error {
	return &errorResponse{
		Status:       status,
		Message:      message,
		Kind:         kind,
		OptionalData: optionalData,
	}
}

func WriteErrorResponse(ctx context.Context, err error, w http.ResponseWriter) []byte {
	w.Header().Set(ContentTypeKey, string(ContentTypeJSON))

	response := make(map[string]any)

	var (
		errMsg string
		status int
		kind   string
	)

	switch e := err.(type) {
	case *errorResponse:
		errMsg = e.Message
		kind = e.Kind
		if http.StatusText(e.Status) != "" {
			status = e.Status
		} else {
			status = http.StatusInternalServerError
		}
		response["optional_data"] = e.OptionalData
	case sharedErrs.BaseError:
		errMsg = e.Message()
		kind = e.Kind().String()
		status = sharedErrs.GetStatusCode(err)
	default:
		errMsg = e.Error()
		kind = sharedErrs.ErrKindUnknown.String()
		status = sharedErrs.GetStatusCode(err)
	}

	response["message"] = errMsg
	response["status"] = status
	response["kind"] = kind

	correlationID := constant.CorrelationIDFromCtx(ctx)
	if correlationID != "" {
		w.Header().Set(constant.XCorrelationIDKey, correlationID)
		response["reference_id"] = correlationID
	}

	result, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		logger.Error(ctx, "Error marshaling http response: %s", marshalErr)
		return []byte(http.StatusText(http.StatusInternalServerError))
	}

	w.WriteHeader(status)
	if _, writeErr := w.Write(result); writeErr != nil {
		logger.Error(ctx, "Error writing http response: %v", writeErr)
		return []byte(http.StatusText(http.StatusInternalServerError))
	}

	return result
}

func WriteHTTPResponse(ctx context.Context, response *Response, w http.ResponseWriter) []byte {
	setHeaders(ctx, response.Headers, w)
	w.WriteHeader(response.Status)

	httpResponse := make(map[string]any)
	httpResponse["status"] = response.Status

	if response.Data != nil {
		httpResponse["data"] = response.Data
	}

	if response.Message != "" {
		httpResponse["message"] = response.Message
	}

	if response.MetaData != nil {
		for k, v := range response.MetaData {
			httpResponse[k] = v
		}
	}

	result, err := json.Marshal(httpResponse)
	if err != nil {
		logger.Error(ctx, "Error marshaling http response: %s", err)
		return []byte(http.StatusText(http.StatusInternalServerError))
	}

	if _, err = w.Write(result); err != nil {
		logger.Error(ctx, "Error writing http response: %v", err)
		return []byte(http.StatusText(http.StatusInternalServerError))
	}

	return result
}
