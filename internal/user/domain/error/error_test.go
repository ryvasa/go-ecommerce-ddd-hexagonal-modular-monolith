package error

import (
	"testing"
)

func TestNewUserError(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		message  string
		category ErrorCategory
	}{
		{
			name:     "validation error",
			code:     "INVALID_EMAIL",
			message:  "email format is invalid",
			category: CategoryValidation,
		},
		{
			name:     "not found error",
			code:     "USER_NOT_FOUND",
			message:  "user does not exist",
			category: CategoryNotFound,
		},
		{
			name:     "conflict error",
			code:     "EMAIL_EXISTS",
			message:  "email already registered",
			category: CategoryConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewUserError(tt.code, tt.message, tt.category)

			if err.Code != tt.code {
				t.Errorf("NewUserError() Code = %q, want %q", err.Code, tt.code)
			}
			if err.Message != tt.message {
				t.Errorf("NewUserError() Message = %q, want %q", err.Message, tt.message)
			}
			if err.Category != tt.category {
				t.Errorf("NewUserError() Category = %v, want %v", err.Category, tt.category)
			}
			if err.Details == nil {
				t.Error("NewUserError() Details should be initialized (not nil)")
			}
		})
	}
}

func TestUserError_Error(t *testing.T) {
	err := NewUserError("TEST_CODE", "test message", CategoryValidation)

	if got := err.Error(); got != "test message" {
		t.Errorf("UserError.Error() = %q, want %q", got, "test message")
	}
}

func TestUserError_GetCode(t *testing.T) {
	err := NewUserError("MY_CODE", "message", CategoryBusiness)

	if got := err.GetCode(); got != "MY_CODE" {
		t.Errorf("UserError.GetCode() = %q, want %q", got, "MY_CODE")
	}
}

func TestUserError_GetCategory(t *testing.T) {
	tests := []struct {
		name     string
		category ErrorCategory
	}{
		{"validation", CategoryValidation},
		{"not found", CategoryNotFound},
		{"conflict", CategoryConflict},
		{"unauthorized", CategoryUnauthorized},
		{"forbidden", CategoryForbidden},
		{"business", CategoryBusiness},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewUserError("CODE", "msg", tt.category)
			if got := err.GetCategory(); got != tt.category {
				t.Errorf("UserError.GetCategory() = %v, want %v", got, tt.category)
			}
		})
	}
}

func TestUserError_WithDetails(t *testing.T) {
	err := NewUserError("CODE", "message", CategoryValidation)
	details := map[string]interface{}{
		"field":  "email",
		"reason": "invalid format",
	}

	result := err.WithDetails(details)

	// Should return same error for chaining
	if result != err {
		t.Error("WithDetails() should return the same error for chaining")
	}

	// Details should be set
	if err.Details["field"] != "email" {
		t.Errorf("Details[field] = %v, want %v", err.Details["field"], "email")
	}
	if err.Details["reason"] != "invalid format" {
		t.Errorf("Details[reason] = %v, want %v", err.Details["reason"], "invalid format")
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("INVALID_INPUT", "input is invalid")

	if err.Category != CategoryValidation {
		t.Errorf("NewValidationError() Category = %v, want %v", err.Category, CategoryValidation)
	}
	if err.Code != "INVALID_INPUT" {
		t.Errorf("NewValidationError() Code = %q, want %q", err.Code, "INVALID_INPUT")
	}
	if err.Message != "input is invalid" {
		t.Errorf("NewValidationError() Message = %q, want %q", err.Message, "input is invalid")
	}
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("USER_NOT_FOUND", "user not found")

	if err.Category != CategoryNotFound {
		t.Errorf("NewNotFoundError() Category = %v, want %v", err.Category, CategoryNotFound)
	}
}

func TestNewConflictError(t *testing.T) {
	err := NewConflictError("DUPLICATE", "resource already exists")

	if err.Category != CategoryConflict {
		t.Errorf("NewConflictError() Category = %v, want %v", err.Category, CategoryConflict)
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	err := NewUnauthorizedError("INVALID_TOKEN", "token is invalid")

	if err.Category != CategoryUnauthorized {
		t.Errorf("NewUnauthorizedError() Category = %v, want %v", err.Category, CategoryUnauthorized)
	}
}

func TestNewForbiddenError(t *testing.T) {
	err := NewForbiddenError("ACCESS_DENIED", "access denied")

	if err.Category != CategoryForbidden {
		t.Errorf("NewForbiddenError() Category = %v, want %v", err.Category, CategoryForbidden)
	}
}

func TestNewBusinessRuleError(t *testing.T) {
	err := NewBusinessRuleError("CART_EMPTY", "cart cannot be empty")

	if err.Category != CategoryBusiness {
		t.Errorf("NewBusinessRuleError() Category = %v, want %v", err.Category, CategoryBusiness)
	}
}
