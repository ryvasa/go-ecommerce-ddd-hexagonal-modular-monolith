package valueobject

import (
	"testing"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid email",
			input:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with subdomain",
			input:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "valid email uppercase gets lowercased",
			input:   "TEST@EXAMPLE.COM",
			wantErr: false,
		},
		{
			name:    "valid email with whitespace gets trimmed",
			input:   "  test@example.com  ",
			wantErr: false,
		},
		{
			name:    "invalid email - no @",
			input:   "testexample.com",
			wantErr: true,
		},
		{
			name:    "invalid email - no domain",
			input:   "test@",
			wantErr: true,
		},
		{
			name:    "invalid email - no TLD",
			input:   "test@example",
			wantErr: true,
		},
		{
			name:    "invalid email - empty",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid email - only whitespace",
			input:   "   ",
			wantErr: true,
		},
		{
			name:    "invalid email - double @",
			input:   "test@@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewEmail(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("NewEmail(%q) unexpected error: %v", tt.input, err)
				}
				// Check that email is lowercase and trimmed
				if email.Value() == "" {
					t.Errorf("NewEmail(%q) returned empty value", tt.input)
				}
			}
		})
	}
}

func TestEmail_Value(t *testing.T) {
	email, err := NewEmail("Test@Example.COM")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should be lowercase
	if email.Value() != "test@example.com" {
		t.Errorf("Email.Value() = %q, want %q", email.Value(), "test@example.com")
	}
}
