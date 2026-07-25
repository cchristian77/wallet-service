package errors

import (
	"fmt"
	"strings"
)

type BaseError interface {
	Error() string
	Message() string
	Kind() Kind
	Cause() error
}

type baseError struct {
	message string
	kind    Kind
	cause   error
	ctxInfo string
}

func (e baseError) Error() string {
	str := ""
	separator := ", "

	if e.kind != "" {
		str += "Kind: {" + string(e.kind) + "}" + separator
	}

	str += e.getMsgAndCauseErrorString()

	if strings.HasSuffix(str, separator) {
		return str[:len(str)-2]
	}

	return str
}

func (e baseError) Message() string {
	return e.message
}

func (e baseError) Cause() error {
	return e.cause
}

func (e baseError) Kind() Kind {
	if e.kind == "" {
		return ErrKindUnknown
	}

	return e.kind
}

func (e baseError) getMsgAndCauseErrorString() string {
	str := ""
	separator := ", "

	if e.ctxInfo == "" {
		if e.message != "" {
			str += "Msg: {" + e.message + "}" + separator
		}
	} else {
		str += "Ctx: {" + e.ctxInfo + "}" + separator
	}

	if e.cause != nil {
		var errStr string

		if c, ok := e.cause.(*baseError); ok {
			errStr = c.getMsgAndCauseErrorString()
		} else {
			errStr = e.cause.Error()
		}

		str += "Cause: {" + errStr + "}" + separator
	}

	if len(str) < len(separator) {
		return ""
	}

	return str[:len(str)-2]
}

func New(kind Kind, message string, args ...any) BaseError {
	return &baseError{
		message: fmt.Sprintf(message, args...),
		kind:    kind,
	}
}

func NewWithCause(kind Kind, message string, err error) BaseError {
	return &baseError{
		message: message,
		kind:    kind,
		cause:   err,
	}
}

func Wrap(err error, ctxInfo string, args ...any) BaseError {
	ctxInfo = fmt.Sprintf(ctxInfo, args...)

	if f, ok := err.(*baseError); ok {
		return &baseError{
			message: f.Message(),
			kind:    f.Kind(),
			ctxInfo: ctxInfo,
			cause:   err,
		}
	}

	var message string
	if err != nil {
		message = err.Error()
	}

	return &baseError{
		message: message,
		ctxInfo: ctxInfo,
		cause:   err,
		kind:    ErrKindUnknown,
	}
}

func E(message string, args ...any) error {
	return &baseError{
		message: fmt.Sprintf(message, args...),
	}
}
