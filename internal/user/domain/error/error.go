package error

// ErrorCategory represents the type of error
type ErrorCategory string

const (
	CategoryValidation   ErrorCategory = "VALIDATION"
	CategoryNotFound     ErrorCategory = "NOT_FOUND"
	CategoryConflict     ErrorCategory = "CONFLICT"
	CategoryUnauthorized ErrorCategory = "UNAUTHORIZED"
	CategoryForbidden    ErrorCategory = "FORBIDDEN"
	CategoryBusiness     ErrorCategory = "BUSINESS_RULE"
)

// UserError represents a domain error with metadata
type UserError struct {
	Code     string
	Message  string
	Category ErrorCategory
	Details  map[string]interface{}
}

func (e *UserError) Error() string {
	return e.Message
}

// WithDetails adds additional context to the error
func (e *UserError) WithDetails(details map[string]interface{}) *UserError {
	e.Details = details
	return e
}

// GetCode returns the error code
func (e *UserError) GetCode() string {
	return e.Code
}

// GetCategory returns the error category
func (e *UserError) GetCategory() ErrorCategory {
	return e.Category
}

// NewUserError creates a new user domain error
func NewUserError(code, message string, category ErrorCategory) *UserError {
	return &UserError{
		Code:     code,
		Message:  message,
		Category: category,
		Details:  make(map[string]interface{}),
	}
}

// Helper functions for common error categories

// NewValidationError creates a validation error
func NewValidationError(code, message string) *UserError {
	return NewUserError(code, message, CategoryValidation)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(code, message string) *UserError {
	return NewUserError(code, message, CategoryNotFound)
}

// NewConflictError creates a conflict error
func NewConflictError(code, message string) *UserError {
	return NewUserError(code, message, CategoryConflict)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(code, message string) *UserError {
	return NewUserError(code, message, CategoryUnauthorized)
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(code, message string) *UserError {
	return NewUserError(code, message, CategoryForbidden)
}

// NewBusinessRuleError creates a business rule error
func NewBusinessRuleError(code, message string) *UserError {
	return NewUserError(code, message, CategoryBusiness)
}
