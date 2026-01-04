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

// CartError represents a domain error with metadata
type CartError struct {
	Code     string
	Message  string
	Category ErrorCategory
	Details  map[string]interface{}
}

func (e *CartError) Error() string {
	return e.Message
}

// WithDetails adds additional context to the error
func (e *CartError) WithDetails(details map[string]interface{}) *CartError {
	e.Details = details
	return e
}

// GetCode returns the error code
func (e *CartError) GetCode() string {
	return e.Code
}

// GetCategory returns the error category
func (e *CartError) GetCategory() ErrorCategory {
	return e.Category
}

// NewCartError creates a new cart domain error
func NewCartError(code, message string, category ErrorCategory) *CartError {
	return &CartError{
		Code:     code,
		Message:  message,
		Category: category,
		Details:  make(map[string]interface{}),
	}
}

// Helper functions for common error categories

// NewValidationError creates a validation error
func NewValidationError(code, message string) *CartError {
	return NewCartError(code, message, CategoryValidation)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(code, message string) *CartError {
	return NewCartError(code, message, CategoryNotFound)
}

// NewConflictError creates a conflict error
func NewConflictError(code, message string) *CartError {
	return NewCartError(code, message, CategoryConflict)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(code, message string) *CartError {
	return NewCartError(code, message, CategoryUnauthorized)
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(code, message string) *CartError {
	return NewCartError(code, message, CategoryForbidden)
}

// NewBusinessRuleError creates a business rule error
func NewBusinessRuleError(code, message string) *CartError {
	return NewCartError(code, message, CategoryBusiness)
}
