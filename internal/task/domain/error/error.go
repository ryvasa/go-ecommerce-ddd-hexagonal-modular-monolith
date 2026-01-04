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

// TaskError represents a domain error with metadata
type TaskError struct {
	Code     string
	Message  string
	Category ErrorCategory
	Details  map[string]interface{}
}

func (e *TaskError) Error() string {
	return e.Message
}

// WithDetails adds additional context to the error
func (e *TaskError) WithDetails(details map[string]interface{}) *TaskError {
	e.Details = details
	return e
}

// GetCode returns the error code
func (e *TaskError) GetCode() string {
	return e.Code
}

// GetCategory returns the error category
func (e *TaskError) GetCategory() ErrorCategory {
	return e.Category
}

// NewTaskError creates a new task domain error
func NewTaskError(code, message string, category ErrorCategory) *TaskError {
	return &TaskError{
		Code:     code,
		Message:  message,
		Category: category,
		Details:  make(map[string]interface{}),
	}
}

// Helper functions for common error categories

// NewValidationError creates a validation error
func NewValidationError(code, message string) *TaskError {
	return NewTaskError(code, message, CategoryValidation)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(code, message string) *TaskError {
	return NewTaskError(code, message, CategoryNotFound)
}

// NewConflictError creates a conflict error
func NewConflictError(code, message string) *TaskError {
	return NewTaskError(code, message, CategoryConflict)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(code, message string) *TaskError {
	return NewTaskError(code, message, CategoryUnauthorized)
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(code, message string) *TaskError {
	return NewTaskError(code, message, CategoryForbidden)
}

// NewBusinessRuleError creates a business rule error
func NewBusinessRuleError(code, message string) *TaskError {
	return NewTaskError(code, message, CategoryBusiness)
}
