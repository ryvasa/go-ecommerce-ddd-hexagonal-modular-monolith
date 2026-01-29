package apperror

import (
	"errors"
	"testing"
)

func TestNotFound(t *testing.T) {
	msg := "user not found"
	err := NotFound(msg)

	if err.Type != ErrorTypeNotFound {
		t.Errorf("NotFound() Type = %v, want %v", err.Type, ErrorTypeNotFound)
	}
	if err.Code != string(ErrorTypeNotFound) {
		t.Errorf("NotFound() Code = %v, want %v", err.Code, string(ErrorTypeNotFound))
	}
	if err.Message != msg {
		t.Errorf("NotFound() Message = %v, want %v", err.Message, msg)
	}
}

func TestValidation(t *testing.T) {
	msg := "email is invalid"
	err := Validation(msg)

	if err.Type != ErrorTypeValidation {
		t.Errorf("Validation() Type = %v, want %v", err.Type, ErrorTypeValidation)
	}
	if err.Code != string(ErrorTypeValidation) {
		t.Errorf("Validation() Code = %v, want %v", err.Code, string(ErrorTypeValidation))
	}
	if err.Message != msg {
		t.Errorf("Validation() Message = %v, want %v", err.Message, msg)
	}
}

func TestInternal(t *testing.T) {
	msg := "database connection failed"
	err := Internal(msg)

	if err.Type != ErrorTypeInternal {
		t.Errorf("Internal() Type = %v, want %v", err.Type, ErrorTypeInternal)
	}
	if err.Code != string(ErrorTypeInternal) {
		t.Errorf("Internal() Code = %v, want %v", err.Code, string(ErrorTypeInternal))
	}
	if err.Message != msg {
		t.Errorf("Internal() Message = %v, want %v", err.Message, msg)
	}
}

func TestConflict(t *testing.T) {
	msg := "email already exists"
	err := Conflict(msg)

	if err.Type != ErrorTypeConflict {
		t.Errorf("Conflict() Type = %v, want %v", err.Type, ErrorTypeConflict)
	}
	if err.Code != string(ErrorTypeConflict) {
		t.Errorf("Conflict() Code = %v, want %v", err.Code, string(ErrorTypeConflict))
	}
	if err.Message != msg {
		t.Errorf("Conflict() Message = %v, want %v", err.Message, msg)
	}
}

func TestUnauthorized(t *testing.T) {
	msg := "invalid token"
	err := Unauthorized(msg)

	if err.Type != ErrorTypeUnauthorized {
		t.Errorf("Unauthorized() Type = %v, want %v", err.Type, ErrorTypeUnauthorized)
	}
	if err.Code != string(ErrorTypeUnauthorized) {
		t.Errorf("Unauthorized() Code = %v, want %v", err.Code, string(ErrorTypeUnauthorized))
	}
	if err.Message != msg {
		t.Errorf("Unauthorized() Message = %v, want %v", err.Message, msg)
	}
}

func TestForbidden(t *testing.T) {
	msg := "access denied"
	err := Forbidden(msg)

	if err.Type != ErrorTypeForbidden {
		t.Errorf("Forbidden() Type = %v, want %v", err.Type, ErrorTypeForbidden)
	}
	if err.Code != string(ErrorTypeForbidden) {
		t.Errorf("Forbidden() Code = %v, want %v", err.Code, string(ErrorTypeForbidden))
	}
	if err.Message != msg {
		t.Errorf("Forbidden() Message = %v, want %v", err.Message, msg)
	}
}

func TestBadRequest(t *testing.T) {
	msg := "invalid request body"
	err := BadRequest(msg)

	if err.Type != ErrorTypeBadRequest {
		t.Errorf("BadRequest() Type = %v, want %v", err.Type, ErrorTypeBadRequest)
	}
	if err.Code != string(ErrorTypeBadRequest) {
		t.Errorf("BadRequest() Code = %v, want %v", err.Code, string(ErrorTypeBadRequest))
	}
	if err.Message != msg {
		t.Errorf("BadRequest() Message = %v, want %v", err.Message, msg)
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		msg         string
		wantType    ErrorType
		wantMessage string
	}{
		{
			name: "wrap app error preserves type",
			err: &Error{
				Type:    ErrorTypeNotFound,
				Code:    "NOT_FOUND",
				Message: "user not found",
			},
			msg:         "failed to get user",
			wantType:    ErrorTypeNotFound,
			wantMessage: "failed to get user: user not found",
		},
		{
			name: "wrap app error with details",
			err: &Error{
				Type:    ErrorTypeValidation,
				Code:    "VALIDATION",
				Message: "invalid email",
				Details: map[string]string{"field": "email"},
			},
			msg:         "registration failed",
			wantType:    ErrorTypeValidation,
			wantMessage: "registration failed: invalid email",
		},
		{
			name:        "wrap standard error becomes internal",
			err:         errors.New("database error"),
			msg:         "failed to save",
			wantType:    ErrorTypeInternal,
			wantMessage: "failed to save: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Wrap(tt.err, tt.msg)

			if result.Type != tt.wantType {
				t.Errorf("Wrap() Type = %v, want %v", result.Type, tt.wantType)
			}
			if result.Message != tt.wantMessage {
				t.Errorf("Wrap() Message = %q, want %q", result.Message, tt.wantMessage)
			}
		})
	}
}

func TestWrap_PreservesDetails(t *testing.T) {
	originalDetails := map[string]string{"field": "email"}
	originalErr := &Error{
		Type:    ErrorTypeValidation,
		Code:    "VALIDATION",
		Message: "invalid",
		Details: originalDetails,
	}

	wrapped := Wrap(originalErr, "wrapped")

	if wrapped.Details == nil {
		t.Error("Wrap() should preserve Details")
	}

	details, ok := wrapped.Details.(map[string]string)
	if !ok {
		t.Fatalf("Details should be map[string]string, got %T", wrapped.Details)
	}

	if details["field"] != "email" {
		t.Errorf("Details[field] = %q, want %q", details["field"], "email")
	}
}
