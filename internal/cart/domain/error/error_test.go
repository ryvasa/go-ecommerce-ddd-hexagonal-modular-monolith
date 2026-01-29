package error

import (
	"testing"
)

func TestNewCartError(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		message  string
		category ErrorCategory
	}{
		{
			name:     "validation error",
			code:     "INVALID_TITLE",
			message:  "title is invalid",
			category: CategoryValidation,
		},
		{
			name:     "not found error",
			code:     "CART_NOT_FOUND",
			message:  "cart does not exist",
			category: CategoryNotFound,
		},
		{
			name:     "conflict error",
			code:     "CART_EXISTS",
			message:  "cart already exists",
			category: CategoryConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewCartError(tt.code, tt.message, tt.category)

			if err.Code != tt.code {
				t.Errorf("NewCartError() Code = %q, want %q", err.Code, tt.code)
			}
			if err.Message != tt.message {
				t.Errorf("NewCartError() Message = %q, want %q", err.Message, tt.message)
			}
			if err.Category != tt.category {
				t.Errorf("NewCartError() Category = %v, want %v", err.Category, tt.category)
			}
			if err.Details == nil {
				t.Error("NewCartError() Details should be initialized (not nil)")
			}
		})
	}
}

func TestCartError_Error(t *testing.T) {
	err := NewCartError("TEST_CODE", "test message", CategoryValidation)

	if got := err.Error(); got != "test message" {
		t.Errorf("CartError.Error() = %q, want %q", got, "test message")
	}
}

func TestCartError_GetCode(t *testing.T) {
	err := NewCartError("MY_CODE", "message", CategoryBusiness)

	if got := err.GetCode(); got != "MY_CODE" {
		t.Errorf("CartError.GetCode() = %q, want %q", got, "MY_CODE")
	}
}

func TestCartError_GetCategory(t *testing.T) {
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
			err := NewCartError("CODE", "msg", tt.category)
			if got := err.GetCategory(); got != tt.category {
				t.Errorf("CartError.GetCategory() = %v, want %v", got, tt.category)
			}
		})
	}
}

func TestCartError_WithDetails(t *testing.T) {
	err := NewCartError("CODE", "message", CategoryValidation)
	details := map[string]interface{}{
		"field":  "title",
		"reason": "too short",
	}

	result := err.WithDetails(details)

	// Should return same error for chaining
	if result != err {
		t.Error("WithDetails() should return the same error for chaining")
	}

	// Details should be set
	if err.Details["field"] != "title" {
		t.Errorf("Details[field] = %v, want %v", err.Details["field"], "title")
	}
	if err.Details["reason"] != "too short" {
		t.Errorf("Details[reason] = %v, want %v", err.Details["reason"], "too short")
	}
}

func TestNewValidationError_Cart(t *testing.T) {
	err := NewValidationError("INVALID_INPUT", "input is invalid")

	if err.Category != CategoryValidation {
		t.Errorf("NewValidationError() Category = %v, want %v", err.Category, CategoryValidation)
	}
	if err.Code != "INVALID_INPUT" {
		t.Errorf("NewValidationError() Code = %q, want %q", err.Code, "INVALID_INPUT")
	}
}

func TestNewNotFoundError_Cart(t *testing.T) {
	err := NewNotFoundError("CART_NOT_FOUND", "cart not found")

	if err.Category != CategoryNotFound {
		t.Errorf("NewNotFoundError() Category = %v, want %v", err.Category, CategoryNotFound)
	}
}

func TestNewConflictError_Cart(t *testing.T) {
	err := NewConflictError("DUPLICATE", "resource already exists")

	if err.Category != CategoryConflict {
		t.Errorf("NewConflictError() Category = %v, want %v", err.Category, CategoryConflict)
	}
}

func TestNewUnauthorizedError_Cart(t *testing.T) {
	err := NewUnauthorizedError("INVALID_TOKEN", "token is invalid")

	if err.Category != CategoryUnauthorized {
		t.Errorf("NewUnauthorizedError() Category = %v, want %v", err.Category, CategoryUnauthorized)
	}
}

func TestNewForbiddenError_Cart(t *testing.T) {
	err := NewForbiddenError("ACCESS_DENIED", "access denied")

	if err.Category != CategoryForbidden {
		t.Errorf("NewForbiddenError() Category = %v, want %v", err.Category, CategoryForbidden)
	}
}

func TestNewBusinessRuleError_Cart(t *testing.T) {
	err := NewBusinessRuleError("CART_EMPTY", "cart cannot be empty")

	if err.Category != CategoryBusiness {
		t.Errorf("NewBusinessRuleError() Category = %v, want %v", err.Category, CategoryBusiness)
	}
}

// Test pre-defined errors
func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      *CartError
		wantCode string
		wantCat  ErrorCategory
	}{
		{
			name:     "ErrCartNotFound",
			err:      ErrCartNotFound,
			wantCode: "CART_NOT_FOUND",
			wantCat:  CategoryNotFound,
		},
		{
			name:     "ErrCartAlreadyExists",
			err:      ErrCartAlreadyExists,
			wantCode: "CART_ALREADY_EXISTS",
			wantCat:  CategoryConflict,
		},
		{
			name:     "ErrEmptyTitle",
			err:      ErrEmptyTitle,
			wantCode: "CART_EMPTY_TITLE",
			wantCat:  CategoryValidation,
		},
		{
			name:     "ErrInvalidCartID",
			err:      ErrInvalidCartID,
			wantCode: "CART_INVALID_ID",
			wantCat:  CategoryValidation,
		},
		{
			name:     "ErrEmptyUserID",
			err:      ErrEmptyUserID,
			wantCode: "CART_EMPTY_USER_ID",
			wantCat:  CategoryValidation,
		},
		{
			name:     "ErrCartAccessDenied",
			err:      ErrCartAccessDenied,
			wantCode: "CART_ACCESS_DENIED",
			wantCat:  CategoryForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.wantCode {
				t.Errorf("%s.Code = %q, want %q", tt.name, tt.err.Code, tt.wantCode)
			}
			if tt.err.Category != tt.wantCat {
				t.Errorf("%s.Category = %v, want %v", tt.name, tt.err.Category, tt.wantCat)
			}
		})
	}
}
