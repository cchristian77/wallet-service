package errors

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

// HTTP ERROR
var (
	// InternalServerErr will throw if any the Internal Server Error happen
	InternalServerErr = New(ErrKindHTTP, "Internal Server Error")

	// NotFoundErr will throw if the requested item is not exists
	NotFoundErr = New(ErrKindDataNotFound, "Requested data is not found")

	// ConflictErr will throw if the current action already exists
	ConflictErr = New(ErrKindRepository, "Requested data already exist")

	// BadParamInputErr will throw if the given request-body or params is not valid
	BadParamInputErr = New(ErrKindBusinessValidation, "Requested parameters are not valid")
)

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = NotFoundErr
	}

	if baseErr, ok := err.(*baseError); ok {
		return getStatusCode(baseErr.kind)
	}

	var be BaseError
	if errors.As(err, &be) {
		return getStatusCode(be.Kind())
	}

	return http.StatusInternalServerError
}

func getStatusCode(k Kind) int {
	switch k {
	case ErrKindValidation, ErrKindBusinessValidation, ErrKindApplication, ErrKindRepository:
		return http.StatusBadRequest
	case ErrKindDataNotFound:
		return http.StatusNotFound
	case ErrKindDatabase, ErrKindHTTP:
		return http.StatusInternalServerError
	case ErrKindInvalidRequest:
		return http.StatusUnprocessableEntity
	case ErrKindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
