package security

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/config"
)

func TestJWTGenerator_Generate(t *testing.T) {
	cfg := config.JWTConfig{
		Secret:     "test-secret-key-12345",
		Expiration: time.Hour,
	}
	generator := NewJWTGenerator(cfg)

	tests := []struct {
		name   string
		userID string
		roles  []string
	}{
		{
			name:   "single role",
			userID: "user-123",
			roles:  []string{"user"},
		},
		{
			name:   "multiple roles",
			userID: "admin-456",
			roles:  []string{"user", "admin"},
		},
		{
			name:   "empty roles",
			userID: "user-789",
			roles:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := generator.Generate(tt.userID, tt.roles)

			if err != nil {
				t.Fatalf("Generate() unexpected error: %v", err)
			}

			if token == "" {
				t.Error("Generate() returned empty token")
			}

			// Token should have 3 parts separated by dots
			parts := strings.Split(token, ".")
			if len(parts) != 3 {
				t.Errorf("Generate() token has %d parts, want 3", len(parts))
			}
		})
	}
}

func TestJWTGenerator_Generate_ContainsCorrectClaims(t *testing.T) {
	cfg := config.JWTConfig{
		Secret:     "test-secret-key-12345",
		Expiration: time.Hour,
	}
	generator := NewJWTGenerator(cfg)

	userID := "user-123"
	roles := []string{"user", "admin"}

	tokenString, err := generator.Generate(userID, roles)
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	// Parse the token to verify claims
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("failed to get claims from token")
	}

	// Check sub claim
	if sub, ok := claims["sub"].(string); !ok || sub != userID {
		t.Errorf("claims[sub] = %v, want %q", claims["sub"], userID)
	}

	// Check roles claim
	rolesInterface, ok := claims["roles"].([]interface{})
	if !ok {
		t.Fatalf("claims[roles] is not []interface{}, got %T", claims["roles"])
	}
	if len(rolesInterface) != len(roles) {
		t.Errorf("claims[roles] has %d items, want %d", len(rolesInterface), len(roles))
	}

	// Check exp claim exists and is in the future
	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal("claims[exp] is not float64")
	}
	expTime := time.Unix(int64(exp), 0)
	if expTime.Before(time.Now()) {
		t.Error("claims[exp] should be in the future")
	}

	// Check iat claim exists
	iat, ok := claims["iat"].(float64)
	if !ok {
		t.Fatal("claims[iat] is not float64")
	}
	iatTime := time.Unix(int64(iat), 0)
	if iatTime.After(time.Now().Add(time.Second)) {
		t.Error("claims[iat] should be now or in the past")
	}
}

func TestJWTGenerator_Generate_TokenIsSignedCorrectly(t *testing.T) {
	cfg := config.JWTConfig{
		Secret:     "correct-secret",
		Expiration: time.Hour,
	}
	generator := NewJWTGenerator(cfg)

	tokenString, err := generator.Generate("user-123", []string{"user"})
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	// Should validate with correct secret
	_, err = jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("correct-secret"), nil
	})
	if err != nil {
		t.Errorf("token should validate with correct secret: %v", err)
	}

	// Should NOT validate with wrong secret
	_, err = jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("wrong-secret"), nil
	})
	if err == nil {
		t.Error("token should NOT validate with wrong secret")
	}
}

func TestJWTGenerator_Generate_ExpirationIsCorrect(t *testing.T) {
	expiration := 2 * time.Hour
	cfg := config.JWTConfig{
		Secret:     "test-secret",
		Expiration: expiration,
	}
	generator := NewJWTGenerator(cfg)

	beforeGenerate := time.Now()
	tokenString, err := generator.Generate("user-123", []string{"user"})
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	exp := time.Unix(int64(claims["exp"].(float64)), 0)

	expectedExp := beforeGenerate.Add(expiration)

	// Allow 1 second tolerance
	if exp.Before(expectedExp.Add(-time.Second)) || exp.After(expectedExp.Add(time.Second)) {
		t.Errorf("exp = %v, want approximately %v", exp, expectedExp)
	}
}

func TestNewJWTGenerator(t *testing.T) {
	cfg := config.JWTConfig{
		Secret:     "test-secret",
		Expiration: time.Hour,
	}

	generator := NewJWTGenerator(cfg)

	if generator == nil {
		t.Error("NewJWTGenerator() returned nil")
	}
}
