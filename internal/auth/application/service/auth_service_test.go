package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth/application/port/in"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/logger"
)

// Mock implementations for testing

type MockUserAuthenticator struct {
	verifyFunc func(ctx context.Context, email, password string) (string, error)
}

func (m *MockUserAuthenticator) VerifyCredentials(ctx context.Context, email, password string) (string, error) {
	if m.verifyFunc != nil {
		return m.verifyFunc(ctx, email, password)
	}
	return "", nil
}

type MockTokenGenerator struct {
	generateFunc func(userID string, roles []string) (string, error)
}

func (m *MockTokenGenerator) Generate(userID string, roles []string) (string, error) {
	if m.generateFunc != nil {
		return m.generateFunc(userID, roles)
	}
	return "mock_token", nil
}

type MockLogger struct{}

func (m *MockLogger) Info(msg string, fields ...any)  {}
func (m *MockLogger) Error(msg string, fields ...any) {}
func (m *MockLogger) With(args ...any) logger.Logger  { return m }

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name          string
		cmd           in.LoginCommand
		userAuth      *MockUserAuthenticator
		tokenGen      *MockTokenGenerator
		wantToken     string
		wantErr       bool
		errContains   string
	}{
		{
			name: "successful login",
			cmd: in.LoginCommand{
				Email:    "test@example.com",
				Password: "password123",
			},
			userAuth: &MockUserAuthenticator{
				verifyFunc: func(ctx context.Context, email, password string) (string, error) {
					return "user-123", nil
				},
			},
			tokenGen: &MockTokenGenerator{
				generateFunc: func(userID string, roles []string) (string, error) {
					if userID != "user-123" {
						t.Errorf("expected userID 'user-123', got %q", userID)
					}
					return "jwt_token_here", nil
				},
			},
			wantToken: "jwt_token_here",
			wantErr:   false,
		},
		{
			name: "invalid credentials - user not found",
			cmd: in.LoginCommand{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			userAuth: &MockUserAuthenticator{
				verifyFunc: func(ctx context.Context, email, password string) (string, error) {
					return "", errors.New("invalid credentials")
				},
			},
			tokenGen:    &MockTokenGenerator{},
			wantToken:   "",
			wantErr:     true,
			errContains: "invalid credentials",
		},
		{
			name: "invalid credentials - wrong password",
			cmd: in.LoginCommand{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			userAuth: &MockUserAuthenticator{
				verifyFunc: func(ctx context.Context, email, password string) (string, error) {
					return "", errors.New("invalid credentials")
				},
			},
			tokenGen:    &MockTokenGenerator{},
			wantToken:   "",
			wantErr:     true,
			errContains: "invalid credentials",
		},
		{
			name: "token generation failure",
			cmd: in.LoginCommand{
				Email:    "test@example.com",
				Password: "password123",
			},
			userAuth: &MockUserAuthenticator{
				verifyFunc: func(ctx context.Context, email, password string) (string, error) {
					return "user-123", nil
				},
			},
			tokenGen: &MockTokenGenerator{
				generateFunc: func(userID string, roles []string) (string, error) {
					return "", errors.New("token generation failed")
				},
			},
			wantToken:   "",
			wantErr:     true,
			errContains: "token generation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &MockLogger{}
			service := NewAuthService(tt.userAuth, tt.tokenGen, logger)

			token, err := service.Login(context.Background(), tt.cmd)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Login() expected error, got nil")
				} else if tt.errContains != "" && !containsSubstring(err.Error(), tt.errContains) {
					t.Errorf("Login() error = %q, want containing %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("Login() unexpected error: %v", err)
				}
				if token != tt.wantToken {
					t.Errorf("Login() token = %q, want %q", token, tt.wantToken)
				}
			}
		})
	}
}

func TestAuthService_Login_VerifiesRolesPassedToTokenGenerator(t *testing.T) {
	var capturedRoles []string

	userAuth := &MockUserAuthenticator{
		verifyFunc: func(ctx context.Context, email, password string) (string, error) {
			return "user-123", nil
		},
	}

	tokenGen := &MockTokenGenerator{
		generateFunc: func(userID string, roles []string) (string, error) {
			capturedRoles = roles
			return "token", nil
		},
	}

	logger := &MockLogger{}
	service := NewAuthService(userAuth, tokenGen, logger)

	cmd := in.LoginCommand{Email: "test@example.com", Password: "password123"}
	_, err := service.Login(context.Background(), cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Currently the service hardcodes "user" role
	if len(capturedRoles) != 1 || capturedRoles[0] != "user" {
		t.Errorf("expected roles [user], got %v", capturedRoles)
	}
}

// Helper function for substring check
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstringHelper(s, substr))
}

func containsSubstringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
