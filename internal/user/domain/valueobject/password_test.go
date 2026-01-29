package valueobject

import (
	"testing"
)

// MockPasswordHasher implements PasswordHasher for testing
type MockPasswordHasher struct {
	hashResult    string
	hashErr       error
	compareResult bool
}

func (m *MockPasswordHasher) Hash(plain string) (string, error) {
	if m.hashErr != nil {
		return "", m.hashErr
	}
	if m.hashResult != "" {
		return m.hashResult, nil
	}
	return "hashed_" + plain, nil
}

func (m *MockPasswordHasher) Compare(hash, plain string) bool {
	return m.compareResult
}

func TestNewPasswordFromPlain(t *testing.T) {
	tests := []struct {
		name    string
		plain   string
		wantErr bool
	}{
		{
			name:    "valid password - 8 characters",
			plain:   "12345678",
			wantErr: false,
		},
		{
			name:    "valid password - longer",
			plain:   "mysecurepassword123",
			wantErr: false,
		},
		{
			name:    "invalid password - too short",
			plain:   "1234567",
			wantErr: true,
		},
		{
			name:    "invalid password - empty",
			plain:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := &MockPasswordHasher{}
			password, err := NewPasswordFromPlain(tt.plain, hasher)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewPasswordFromPlain(%q) expected error, got nil", tt.plain)
				}
			} else {
				if err != nil {
					t.Errorf("NewPasswordFromPlain(%q) unexpected error: %v", tt.plain, err)
				}
				if password.Hash() == "" {
					t.Errorf("NewPasswordFromPlain(%q) returned empty hash", tt.plain)
				}
			}
		})
	}
}

func TestNewHashedPassword(t *testing.T) {
	tests := []struct {
		name    string
		hash    string
		wantErr bool
	}{
		{
			name:    "valid hash",
			hash:    "$2a$10$somehash",
			wantErr: false,
		},
		{
			name:    "empty hash",
			hash:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := NewHashedPassword(tt.hash)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewHashedPassword(%q) expected error, got nil", tt.hash)
				}
			} else {
				if err != nil {
					t.Errorf("NewHashedPassword(%q) unexpected error: %v", tt.hash, err)
				}
				if password.Hash() != tt.hash {
					t.Errorf("Password.Hash() = %q, want %q", password.Hash(), tt.hash)
				}
			}
		})
	}
}

func TestPassword_Verify(t *testing.T) {
	tests := []struct {
		name          string
		compareResult bool
		want          bool
	}{
		{
			name:          "correct password",
			compareResult: true,
			want:          true,
		},
		{
			name:          "wrong password",
			compareResult: false,
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := &MockPasswordHasher{compareResult: tt.compareResult}
			password, _ := NewHashedPassword("somehash")

			got := password.Verify("plain", hasher)
			if got != tt.want {
				t.Errorf("Password.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}
