package apperror

// ErrorType represents the category of error
type ErrorType string

const (
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
	ErrorTypeConflict     ErrorType = "CONFLICT"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
	ErrorTypeBadRequest   ErrorType = "BAD_REQUEST"
	ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"
)

type Error struct {
	Type    ErrorType   `json:"-"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

// IsNotFound checks if the error is a NotFound type
func (e *Error) IsNotFound() bool {
	return e.Type == ErrorTypeNotFound
}

// IsValidation checks if the error is a Validation type
func (e *Error) IsValidation() bool {
	return e.Type == ErrorTypeValidation
}

// IsConflict checks if the error is a Conflict type
func (e *Error) IsConflict() bool {
	return e.Type == ErrorTypeConflict
}

// IsUnauthorized checks if the error is an Unauthorized type
func (e *Error) IsUnauthorized() bool {
	return e.Type == ErrorTypeUnauthorized
}

// IsForbidden checks if the error is a Forbidden type
func (e *Error) IsForbidden() bool {
	return e.Type == ErrorTypeForbidden
}

// IsBadRequest checks if the error is a BadRequest type
func (e *Error) IsBadRequest() bool {
	return e.Type == ErrorTypeBadRequest
}

// WithDetails adds additional details to the error
func (e *Error) WithDetails(details interface{}) *Error {
	e.Details = details
	return e
}
