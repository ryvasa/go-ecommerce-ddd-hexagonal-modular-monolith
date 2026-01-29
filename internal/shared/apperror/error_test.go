package apperror

import (
	"testing"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "returns message",
			err: &Error{
				Type:    ErrorTypeNotFound,
				Code:    "NOT_FOUND",
				Message: "resource not found",
			},
			want: "resource not found",
		},
		{
			name: "empty message",
			err: &Error{
				Type:    ErrorTypeInternal,
				Code:    "INTERNAL",
				Message: "",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error.Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestError_IsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want bool
	}{
		{
			name: "is not found",
			err:  &Error{Type: ErrorTypeNotFound},
			want: true,
		},
		{
			name: "is not not found - validation",
			err:  &Error{Type: ErrorTypeValidation},
			want: false,
		},
		{
			name: "is not not found - internal",
			err:  &Error{Type: ErrorTypeInternal},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsNotFound(); got != tt.want {
				t.Errorf("Error.IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_IsValidation(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want bool
	}{
		{
			name: "is validation",
			err:  &Error{Type: ErrorTypeValidation},
			want: true,
		},
		{
			name: "is not validation - not found",
			err:  &Error{Type: ErrorTypeNotFound},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsValidation(); got != tt.want {
				t.Errorf("Error.IsValidation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_IsConflict(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want bool
	}{
		{
			name: "is conflict",
			err:  &Error{Type: ErrorTypeConflict},
			want: true,
		},
		{
			name: "is not conflict",
			err:  &Error{Type: ErrorTypeNotFound},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsConflict(); got != tt.want {
				t.Errorf("Error.IsConflict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_IsUnauthorized(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want bool
	}{
		{
			name: "is unauthorized",
			err:  &Error{Type: ErrorTypeUnauthorized},
			want: true,
		},
		{
			name: "is not unauthorized",
			err:  &Error{Type: ErrorTypeForbidden},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsUnauthorized(); got != tt.want {
				t.Errorf("Error.IsUnauthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_IsForbidden(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want bool
	}{
		{
			name: "is forbidden",
			err:  &Error{Type: ErrorTypeForbidden},
			want: true,
		},
		{
			name: "is not forbidden",
			err:  &Error{Type: ErrorTypeUnauthorized},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsForbidden(); got != tt.want {
				t.Errorf("Error.IsForbidden() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_IsBadRequest(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want bool
	}{
		{
			name: "is bad request",
			err:  &Error{Type: ErrorTypeBadRequest},
			want: true,
		},
		{
			name: "is not bad request",
			err:  &Error{Type: ErrorTypeValidation},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsBadRequest(); got != tt.want {
				t.Errorf("Error.IsBadRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_WithDetails(t *testing.T) {
	err := &Error{
		Type:    ErrorTypeValidation,
		Code:    "VALIDATION",
		Message: "validation failed",
	}

	details := map[string]string{"field": "email", "reason": "invalid format"}
	result := err.WithDetails(details)

	// Should return the same error for chaining
	if result != err {
		t.Error("WithDetails() should return the same error for chaining")
	}

	// Details should be set
	if err.Details == nil {
		t.Error("WithDetails() should set Details field")
	}

	detailsMap, ok := err.Details.(map[string]string)
	if !ok {
		t.Fatalf("Details should be map[string]string, got %T", err.Details)
	}

	if detailsMap["field"] != "email" {
		t.Errorf("Details[field] = %q, want %q", detailsMap["field"], "email")
	}
}
