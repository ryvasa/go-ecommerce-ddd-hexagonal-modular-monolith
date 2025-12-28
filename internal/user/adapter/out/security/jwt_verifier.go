package security

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/middleware"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/config"
)

type JWTVerifierImpl struct {
	secret []byte
}

func NewJWTVerifier(cfg config.JWTConfig) middleware.JWTVerifier {
	return &JWTVerifierImpl{
		secret: []byte(cfg.Secret),
	}
}

func (v *JWTVerifierImpl) Verify(authHeader string) (middleware.Claims, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return middleware.Claims{}, errors.New("invalid auth header")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return v.secret, nil
	})
	if err != nil || !token.Valid {
		return middleware.Claims{}, errors.New("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)

	return middleware.Claims{
		UserID: claims["sub"].(string),
		Role:   claims["role"].(string),
	}, nil
}
