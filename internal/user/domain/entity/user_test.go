package entity

import (
	"testing"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
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

func createTestUser(t *testing.T) *User {
	t.Helper()
	email, err := valueobject.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	password, err := valueobject.NewHashedPassword("hashedpassword")
	if err != nil {
		t.Fatalf("failed to create password: %v", err)
	}
	return NewUser(email, password, "testuser", "John", "Doe")
}

func TestNewUser(t *testing.T) {
	email, _ := valueobject.NewEmail("test@example.com")
	password, _ := valueobject.NewHashedPassword("hashedpassword")

	user := NewUser(email, password, "testuser", "John", "Doe")

	if user.ID() == "" {
		t.Error("NewUser() should generate ID")
	}
	if user.Email().Value() != "test@example.com" {
		t.Errorf("User.Email() = %q, want %q", user.Email().Value(), "test@example.com")
	}
	if user.Username() != "testuser" {
		t.Errorf("User.Username() = %q, want %q", user.Username(), "testuser")
	}
	if user.FirstName() != "John" {
		t.Errorf("User.FirstName() = %q, want %q", user.FirstName(), "John")
	}
	if user.LastName() != "Doe" {
		t.Errorf("User.LastName() = %q, want %q", user.LastName(), "Doe")
	}

	// New user should have "user" role by default
	roles := user.Roles()
	hasUserRole := false
	for _, r := range roles {
		if r == "user" {
			hasUserRole = true
			break
		}
	}
	if !hasUserRole {
		t.Error("NewUser() should assign 'user' role by default")
	}
}

func TestUser_AssignRole(t *testing.T) {
	tests := []struct {
		name         string
		existingRole string
		newRole      string
		wantErr      bool
	}{
		{
			name:         "assign new role",
			existingRole: "",
			newRole:      "admin",
			wantErr:      false,
		},
		{
			name:         "assign duplicate role",
			existingRole: "admin",
			newRole:      "admin",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := createTestUser(t)

			// Pre-assign role if needed
			if tt.existingRole != "" && tt.existingRole != "user" {
				_ = user.AssignRole(tt.existingRole)
			}

			err := user.AssignRole(tt.newRole)

			if tt.wantErr {
				if err == nil {
					t.Errorf("User.AssignRole(%q) expected error, got nil", tt.newRole)
				}
			} else {
				if err != nil {
					t.Errorf("User.AssignRole(%q) unexpected error: %v", tt.newRole, err)
				}
				// Verify role was assigned
				roles := user.Roles()
				hasRole := false
				for _, r := range roles {
					if r == tt.newRole {
						hasRole = true
						break
					}
				}
				if !hasRole {
					t.Errorf("User.Roles() should contain %q", tt.newRole)
				}
			}
		})
	}
}

func TestUser_ChangePassword(t *testing.T) {
	tests := []struct {
		name         string
		oldPlain     string
		newPlain     string
		verifyResult bool
		wantErr      bool
	}{
		{
			name:         "valid password change",
			oldPlain:     "oldpassword",
			newPlain:     "newpassword123",
			verifyResult: true,
			wantErr:      false,
		},
		{
			name:         "wrong current password",
			oldPlain:     "wrongpassword",
			newPlain:     "newpassword123",
			verifyResult: false,
			wantErr:      true,
		},
		{
			name:         "new password too short",
			oldPlain:     "oldpassword",
			newPlain:     "short",
			verifyResult: true,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, _ := valueobject.NewEmail("test@example.com")
			password, _ := valueobject.NewHashedPassword("oldhash")
			user := RehydrateUser("user-id", email, password)

			hasher := &MockPasswordHasher{compareResult: tt.verifyResult}

			err := user.ChangePassword(tt.oldPlain, tt.newPlain, hasher)

			if tt.wantErr {
				if err == nil {
					t.Error("User.ChangePassword() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("User.ChangePassword() unexpected error: %v", err)
				}
				// Password should be changed (hash should be different)
				if user.Password().Hash() == "oldhash" {
					t.Error("User.Password() should be changed after ChangePassword()")
				}
			}
		})
	}
}

func TestRehydrateUser(t *testing.T) {
	email, _ := valueobject.NewEmail("test@example.com")
	password, _ := valueobject.NewHashedPassword("hashedpassword")

	user := RehydrateUser("user-123", email, password)

	if user.ID() != "user-123" {
		t.Errorf("User.ID() = %q, want %q", user.ID(), "user-123")
	}
	if user.Email().Value() != "test@example.com" {
		t.Errorf("User.Email() = %q, want %q", user.Email().Value(), "test@example.com")
	}
}
