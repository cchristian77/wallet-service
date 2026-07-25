package errors

// Kind is used to pinpoint a specific error.
type Kind string

func (k Kind) String() string {
	return string(k)
}

const (
	ErrKindHTTP           Kind = "http_error"
	ErrKindInvalidRequest Kind = "invalid_request_error"
	ErrKindValidation     Kind = "validation_error"

	ErrKindApplication        Kind = "application_error"
	ErrKindBusinessValidation Kind = "business_validation_error"

	ErrKindRepository   Kind = "repository_error"
	ErrKindDatabase     Kind = "database_error"
	ErrKindDataNotFound Kind = "data_not_found_error"

	ErrKindConflict Kind = "conflict_error"
	ErrKindUnknown  Kind = "unknown_error"
)
