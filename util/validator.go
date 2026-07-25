package util

import (
	"errors"
	"net/http"

	sharedErrs "github.com/cchristian77/wallet-service/shared/errors"
	"github.com/cchristian77/wallet-service/shared/fhttp"
	"github.com/go-playground/validator/v10"
)

var v *validator.Validate

func init() {
	v = validator.New()
}

func Validate(input any) error {
	if v == nil {
		v = validator.New()
	}

	err := v.Struct(input)
	if err == nil {
		return nil
	}

	var fieldErrs []fhttp.OptionalData
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, e := range validationErrs {
			fieldErrs = append(fieldErrs, fhttp.OptionalData{
				Key:   e.Field(),
				Value: e.Error(),
			})
		}
	}

	return fhttp.NewErrorResponse(
		http.StatusBadRequest,
		sharedErrs.ErrKindValidation.String(),
		"Validation Error",
		fieldErrs...,
	)
}
