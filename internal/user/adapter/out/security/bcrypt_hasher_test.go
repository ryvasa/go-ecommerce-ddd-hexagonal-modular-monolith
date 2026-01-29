package security

import (
	"testing"
)

func TestBcryptHasher_Hash(t *testing.T) {
	hasher := NewBcryptHasher()

	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "hash simple password",
			password: "password123",
		},
		{
			name:     "hash complex password",
			password: "MyP@ssw0rd!#$%^&*()",
		},
		{
			name:     "hash long password",
			password: "thisisaverylongpasswordthatexceedstheusuallengthlimit123456789",
		},
		{
			name:     "hash unicode password",
			password: "密碼パスワード🔐",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := hasher.Hash(tt.password)

			if err != nil {
				t.Fatalf("Hash() unexpected error: %v", err)
			}

			if hash == "" {
				t.Error("Hash() returned empty string")
			}

			// Hash should start with bcrypt prefix
			if len(hash) < 4 || hash[:4] != "$2a$" {
				t.Errorf("Hash() = %q, doesn't look like bcrypt hash", hash)
			}

			// Hash should be different from plain password
			if hash == tt.password {
				t.Error("Hash() returned the same as plain password")
			}
		})
	}
}

func TestBcryptHasher_Hash_DifferentHashes(t *testing.T) {
	hasher := NewBcryptHasher()
	password := "samepassword"

	hash1, err1 := hasher.Hash(password)
	hash2, err2 := hasher.Hash(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("Hash() unexpected errors: %v, %v", err1, err2)
	}

	// Same password should produce different hashes (because of salt)
	if hash1 == hash2 {
		t.Error("Hash() should produce different hashes for same password due to salt")
	}
}

func TestBcryptHasher_Compare(t *testing.T) {
	hasher := NewBcryptHasher()

	tests := []struct {
		name       string
		password   string
		checkWith  string
		wantResult bool
	}{
		{
			name:       "correct password",
			password:   "correctpassword",
			checkWith:  "correctpassword",
			wantResult: true,
		},
		{
			name:       "wrong password",
			password:   "correctpassword",
			checkWith:  "wrongpassword",
			wantResult: false,
		},
		{
			name:       "empty password check",
			password:   "somepassword",
			checkWith:  "",
			wantResult: false,
		},
		{
			name:       "case sensitive",
			password:   "Password123",
			checkWith:  "password123",
			wantResult: false,
		},
		{
			name:       "with special characters",
			password:   "p@ssw0rd!",
			checkWith:  "p@ssw0rd!",
			wantResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First hash the password
			hash, err := hasher.Hash(tt.password)
			if err != nil {
				t.Fatalf("Hash() unexpected error: %v", err)
			}

			// Then compare
			result := hasher.Compare(hash, tt.checkWith)

			if result != tt.wantResult {
				t.Errorf("Compare() = %v, want %v", result, tt.wantResult)
			}
		})
	}
}

func TestBcryptHasher_Compare_InvalidHash(t *testing.T) {
	hasher := NewBcryptHasher()

	tests := []struct {
		name string
		hash string
	}{
		{
			name: "empty hash",
			hash: "",
		},
		{
			name: "invalid hash format",
			hash: "not-a-valid-bcrypt-hash",
		},
		{
			name: "plain text",
			hash: "password123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasher.Compare(tt.hash, "password")

			if result {
				t.Error("Compare() should return false for invalid hash")
			}
		})
	}
}

func TestNewBcryptHasher(t *testing.T) {
	hasher := NewBcryptHasher()

	if hasher == nil {
		t.Error("NewBcryptHasher() returned nil")
	}
}
