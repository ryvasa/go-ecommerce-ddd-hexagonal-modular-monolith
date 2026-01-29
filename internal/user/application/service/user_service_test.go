package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/logger"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/in"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
	usererror "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/error"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
)

// Mock implementations for testing

type MockUserRepository struct {
	saveFunc        func(ctx context.Context, user *entity.User) error
	findByEmailFunc func(ctx context.Context, email valueobject.Email) (*entity.User, error)
	existsByIDFunc  func(ctx context.Context, id string) (bool, error)
	getByIDFunc     func(ctx context.Context, id string) (*entity.User, error)
}

func (m *MockUserRepository) Save(ctx context.Context, user *entity.User) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
	if m.findByEmailFunc != nil {
		return m.findByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	if m.existsByIDFunc != nil {
		return m.existsByIDFunc(ctx, id)
	}
	return false, nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

type MockPasswordHasher struct {
	hashFunc    func(plain string) (string, error)
	compareFunc func(hash, plain string) bool
}

func (m *MockPasswordHasher) Hash(plain string) (string, error) {
	if m.hashFunc != nil {
		return m.hashFunc(plain)
	}
	return "hashed_password", nil
}

func (m *MockPasswordHasher) Compare(hash, plain string) bool {
	if m.compareFunc != nil {
		return m.compareFunc(hash, plain)
	}
	return true
}

type MockLogger struct{}

func (m *MockLogger) Info(msg string, fields ...any)  {}
func (m *MockLogger) Error(msg string, fields ...any) {}
func (m *MockLogger) With(args ...any) logger.Logger  { return m }

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name    string
		cmd     in.RegisterUserCommand
		repo    *MockUserRepository
		wantErr error
	}{
		{
			name: "successful registration",
			cmd: in.RegisterUserCommand{
				Email:     "test@example.com",
				Password:  "password123",
				Username:  "testuser",
				FirstName: "John",
				LastName:  "Doe",
			},
			repo: &MockUserRepository{
				findByEmailFunc: func(ctx context.Context, email valueobject.Email) (*entity.User, error) {
					return nil, nil // No existing user
				},
				saveFunc: func(ctx context.Context, user *entity.User) error {
					return nil
				},
			},
			wantErr: nil,
		},
		{
			name: "email already exists",
			cmd: in.RegisterUserCommand{
				Email:     "existing@example.com",
				Password:  "password123",
				Username:  "testuser",
				FirstName: "John",
				LastName:  "Doe",
			},
			repo: &MockUserRepository{
				findByEmailFunc: func(ctx context.Context, email valueobject.Email) (*entity.User, error) {
					// Return existing user
					e, _ := valueobject.NewEmail("existing@example.com")
					p, _ := valueobject.NewHashedPassword("hash")
					return entity.RehydrateUser("existing-id", e, p), nil
				},
			},
			wantErr: usererror.ErrEmailAlreadyUsed,
		},
		{
			name: "invalid email format",
			cmd: in.RegisterUserCommand{
				Email:     "invalid-email",
				Password:  "password123",
				Username:  "testuser",
				FirstName: "John",
				LastName:  "Doe",
			},
			repo:    &MockUserRepository{},
			wantErr: errors.New("invalid email format"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := &MockPasswordHasher{}
			logger := &MockLogger{}
			service := NewUserService(tt.repo, hasher, logger)

			err := service.Register(context.Background(), tt.cmd)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Register() expected error, got nil")
				} else if tt.wantErr.Error() != err.Error() && !errors.Is(err, tt.wantErr) {
					// Allow either exact match or wrapped error
					if tt.wantErr.Error() != err.Error() {
						t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Register() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestUserService_VerifyCredentials(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		password   string
		repo       *MockUserRepository
		hasher     *MockPasswordHasher
		wantUserID string
		wantErr    error
	}{
		{
			name:     "successful verification",
			email:    "test@example.com",
			password: "correctpassword",
			repo: &MockUserRepository{
				findByEmailFunc: func(ctx context.Context, email valueobject.Email) (*entity.User, error) {
					e, _ := valueobject.NewEmail("test@example.com")
					p, _ := valueobject.NewHashedPassword("hashedpassword")
					return entity.RehydrateUser("user-123", e, p), nil
				},
			},
			hasher: &MockPasswordHasher{
				compareFunc: func(hash, plain string) bool { return true },
			},
			wantUserID: "user-123",
			wantErr:    nil,
		},
		{
			name:     "user not found",
			email:    "notfound@example.com",
			password: "password",
			repo: &MockUserRepository{
				findByEmailFunc: func(ctx context.Context, email valueobject.Email) (*entity.User, error) {
					return nil, nil
				},
			},
			hasher:     &MockPasswordHasher{},
			wantUserID: "",
			wantErr:    usererror.ErrInvalidCredential,
		},
		{
			name:     "wrong password",
			email:    "test@example.com",
			password: "wrongpassword",
			repo: &MockUserRepository{
				findByEmailFunc: func(ctx context.Context, email valueobject.Email) (*entity.User, error) {
					e, _ := valueobject.NewEmail("test@example.com")
					p, _ := valueobject.NewHashedPassword("hashedpassword")
					return entity.RehydrateUser("user-123", e, p), nil
				},
			},
			hasher: &MockPasswordHasher{
				compareFunc: func(hash, plain string) bool { return false },
			},
			wantUserID: "",
			wantErr:    usererror.ErrInvalidCredential,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &MockLogger{}
			service := NewUserService(tt.repo, tt.hasher, logger)

			email, _ := valueobject.NewEmail(tt.email)
			userID, err := service.VerifyCredentials(context.Background(), email, tt.password)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("VerifyCredentials() expected error, got nil")
				} else if !errors.Is(err, tt.wantErr) {
					t.Errorf("VerifyCredentials() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("VerifyCredentials() unexpected error: %v", err)
				}
				if userID != tt.wantUserID {
					t.Errorf("VerifyCredentials() userID = %v, want %v", userID, tt.wantUserID)
				}
			}
		})
	}
}

func TestUserService_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		repo    *MockUserRepository
		wantErr error
	}{
		{
			name:   "user found",
			userID: "user-123",
			repo: &MockUserRepository{
				getByIDFunc: func(ctx context.Context, id string) (*entity.User, error) {
					e, _ := valueobject.NewEmail("test@example.com")
					p, _ := valueobject.NewHashedPassword("hash")
					return entity.RehydrateUser("user-123", e, p), nil
				},
			},
			wantErr: nil,
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			repo: &MockUserRepository{
				getByIDFunc: func(ctx context.Context, id string) (*entity.User, error) {
					return nil, nil
				},
			},
			wantErr: usererror.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := &MockPasswordHasher{}
			logger := &MockLogger{}
			service := NewUserService(tt.repo, hasher, logger)

			user, err := service.GetByID(context.Background(), tt.userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetByID() expected error, got nil")
				} else if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("GetByID() unexpected error: %v", err)
				}
				if user == nil {
					t.Error("GetByID() returned nil user")
				}
			}
		})
	}
}
