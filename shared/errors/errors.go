package errors

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func NewBusinessValidationErr(message string, args ...any) BaseError {
	return &baseError{
		message: fmt.Sprintf(message, args...),
		kind:    ErrKindBusinessValidation,
	}
}

func NewRepositoryErr(err error, message string, args ...any) BaseError {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFoundErr
	}

	return NewWithCause(ErrKindRepository, fmt.Sprintf(message, args...), err)
}
