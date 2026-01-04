package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/config"
)

type JWTGenerator struct {
	secret     []byte
	expiration time.Duration
}

func NewJWTGenerator(cfg config.JWTConfig) *JWTGenerator {
	return &JWTGenerator{
		secret:     []byte(cfg.Secret),
		expiration: cfg.Expiration,
	}
}

func (j *JWTGenerator) Generate(userID string, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"roles": roles, // Changed from "role" to "roles" and accepting slice
		"exp":   time.Now().Add(j.expiration).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}
