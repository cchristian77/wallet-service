package errors

// Kind is used to pinpoint a specific error.
type Kind string

func (k Kind) String() string {
	return string(k)
}

const (
	// HTTP or Controller Errors
	ErrKindHTTP           Kind = "http_error"
	ErrKindInvalidRequest Kind = "invalid_request_error"
	ErrKindValidation     Kind = "validation_error"

	// ErrKindApplication errors are errors that might be resolved by retrying the same request at a later time
	ErrKindApplication        Kind = "application_error"
	ErrKindBusinessValidation Kind = "business_validation_error"

	// Repository Error
	ErrKindRepository   Kind = "repository_error"
	ErrKindDatabase     Kind = "database_error"
	ErrKindDataNotFound Kind = "data_not_found_error"

	// Utilities or Miscellaneous Errors
	ErrKindConflict Kind = "conflict_error"
	ErrKindUnknown  Kind = "unknown_error"
)
