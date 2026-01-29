package security

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/config"
)

func createTestToken(secret string, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func TestJWTVerifier_Verify(t *testing.T) {
	secret := "test-secret-key-12345"
	cfg := config.JWTConfig{
		Secret:     secret,
		Expiration: time.Hour,
	}
	verifier := NewJWTVerifier(cfg)

	tests := []struct {
		name        string
		authHeader  string
		wantUserID  string
		wantRole    string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid token",
			authHeader: "Bearer " + createTestToken(secret, jwt.MapClaims{
				"sub":  "user-123",
				"role": "user",
				"exp":  time.Now().Add(time.Hour).Unix(),
			}),
			wantUserID: "user-123",
			wantRole:   "user",
			wantErr:    false,
		},
		{
			name: "valid token with admin role",
			authHeader: "Bearer " + createTestToken(secret, jwt.MapClaims{
				"sub":  "admin-456",
				"role": "admin",
				"exp":  time.Now().Add(time.Hour).Unix(),
			}),
			wantUserID: "admin-456",
			wantRole:   "admin",
			wantErr:    false,
		},
		{
			name:        "missing Bearer prefix",
			authHeader:  createTestToken(secret, jwt.MapClaims{"sub": "user-123", "role": "user", "exp": time.Now().Add(time.Hour).Unix()}),
			wantErr:     true,
			errContains: "invalid auth header",
		},
		{
			name:        "empty auth header",
			authHeader:  "",
			wantErr:     true,
			errContains: "invalid auth header",
		},
		{
			name:        "only Bearer word",
			authHeader:  "Bearer",
			wantErr:     true,
			errContains: "invalid",
		},
		{
			name:        "invalid token format",
			authHeader:  "Bearer invalid-token-string",
			wantErr:     true,
			errContains: "invalid token",
		},
		{
			name: "expired token",
			authHeader: "Bearer " + createTestToken(secret, jwt.MapClaims{
				"sub":  "user-123",
				"role": "user",
				"exp":  time.Now().Add(-time.Hour).Unix(), // expired 1 hour ago
			}),
			wantErr:     true,
			errContains: "invalid token",
		},
		{
			name: "wrong signature",
			authHeader: "Bearer " + createTestToken("wrong-secret", jwt.MapClaims{
				"sub":  "user-123",
				"role": "user",
				"exp":  time.Now().Add(time.Hour).Unix(),
			}),
			wantErr:     true,
			errContains: "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := verifier.Verify(tt.authHeader)

			if tt.wantErr {
				if err == nil {
					t.Error("Verify() expected error, got nil")
				} else if tt.errContains != "" && !containsStr(err.Error(), tt.errContains) {
					t.Errorf("Verify() error = %q, want containing %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("Verify() unexpected error: %v", err)
				}
				if claims.UserID != tt.wantUserID {
					t.Errorf("Claims.UserID = %q, want %q", claims.UserID, tt.wantUserID)
				}
				if claims.Role != tt.wantRole {
					t.Errorf("Claims.Role = %q, want %q", claims.Role, tt.wantRole)
				}
			}
		})
	}
}

func TestNewJWTVerifier(t *testing.T) {
	cfg := config.JWTConfig{
		Secret:     "test-secret",
		Expiration: time.Hour,
	}

	verifier := NewJWTVerifier(cfg)

	if verifier == nil {
		t.Error("NewJWTVerifier() returned nil")
	}
}

// Helper function
func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
