//go:build integration

package persistence

import (
	"context"
	"testing"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Auto migrate the user model
	if err := db.AutoMigrate(&UserModel{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func createTestUserEntity(t *testing.T, email string) *entity.User {
	t.Helper()
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	passwordVO, err := valueobject.NewHashedPassword("$2a$10$hashedpassword")
	if err != nil {
		t.Fatalf("failed to create password: %v", err)
	}
	return entity.NewUser(emailVO, passwordVO, "testuser", "John", "Doe")
}

func TestGormUserRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db)

	tests := []struct {
		name    string
		user    func(t *testing.T) *entity.User
		wantErr bool
	}{
		{
			name: "save new user successfully",
			user: func(t *testing.T) *entity.User {
				return createTestUserEntity(t, "test@example.com")
			},
			wantErr: false,
		},
		{
			name: "save user with different email",
			user: func(t *testing.T) *entity.User {
				return createTestUserEntity(t, "another@example.com")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := tt.user(t)
			err := repo.Save(context.Background(), user)

			if tt.wantErr {
				if err == nil {
					t.Error("Save() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Save() unexpected error: %v", err)
				}

				// Verify user was saved
				var model UserModel
				if err := db.First(&model, "id = ?", user.ID()).Error; err != nil {
					t.Errorf("failed to find saved user: %v", err)
				}
				if model.Email != user.Email().Value() {
					t.Errorf("saved email = %q, want %q", model.Email, user.Email().Value())
				}
			}
		})
	}
}

func TestGormUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db)

	// Setup: save a test user
	existingUser := createTestUserEntity(t, "existing@example.com")
	if err := repo.Save(context.Background(), existingUser); err != nil {
		t.Fatalf("failed to setup test user: %v", err)
	}

	tests := []struct {
		name      string
		email     string
		wantFound bool
		wantErr   bool
	}{
		{
			name:      "find existing user",
			email:     "existing@example.com",
			wantFound: true,
			wantErr:   false,
		},
		{
			name:      "user not found",
			email:     "notfound@example.com",
			wantFound: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emailVO, _ := valueobject.NewEmail(tt.email)
			user, err := repo.FindByEmail(context.Background(), emailVO)

			if tt.wantErr {
				if err == nil {
					t.Error("FindByEmail() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("FindByEmail() unexpected error: %v", err)
				}

				if tt.wantFound {
					if user == nil {
						t.Error("FindByEmail() expected user, got nil")
					} else if user.Email().Value() != tt.email {
						t.Errorf("user.Email() = %q, want %q", user.Email().Value(), tt.email)
					}
				} else {
					if user != nil {
						t.Errorf("FindByEmail() expected nil, got user with email %q", user.Email().Value())
					}
				}
			}
		})
	}
}

func TestGormUserRepository_ExistsByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db)

	// Setup: save a test user
	existingUser := createTestUserEntity(t, "exists@example.com")
	if err := repo.Save(context.Background(), existingUser); err != nil {
		t.Fatalf("failed to setup test user: %v", err)
	}

	tests := []struct {
		name       string
		userID     string
		wantExists bool
		wantErr    bool
	}{
		{
			name:       "existing user",
			userID:     existingUser.ID(),
			wantExists: true,
			wantErr:    false,
		},
		{
			name:       "non-existing user",
			userID:     "non-existing-id",
			wantExists: false,
			wantErr:    false,
		},
		{
			name:       "empty id",
			userID:     "",
			wantExists: false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := repo.ExistsByID(context.Background(), tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Error("ExistsByID() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ExistsByID() unexpected error: %v", err)
				}
				if exists != tt.wantExists {
					t.Errorf("ExistsByID() = %v, want %v", exists, tt.wantExists)
				}
			}
		})
	}
}

func TestGormUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db)

	// Setup: save a test user
	existingUser := createTestUserEntity(t, "getbyid@example.com")
	if err := repo.Save(context.Background(), existingUser); err != nil {
		t.Fatalf("failed to setup test user: %v", err)
	}

	tests := []struct {
		name      string
		userID    string
		wantFound bool
		wantErr   bool
	}{
		{
			name:      "get existing user",
			userID:    existingUser.ID(),
			wantFound: true,
			wantErr:   false,
		},
		{
			name:      "user not found",
			userID:    "non-existing-id",
			wantFound: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByID(context.Background(), tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Error("GetByID() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("GetByID() unexpected error: %v", err)
				}

				if tt.wantFound {
					if user == nil {
						t.Error("GetByID() expected user, got nil")
					} else if user.ID() != tt.userID {
						t.Errorf("user.ID() = %q, want %q", user.ID(), tt.userID)
					}
				} else {
					if user != nil {
						t.Errorf("GetByID() expected nil, got user with id %q", user.ID())
					}
				}
			}
		})
	}
}
