package config

import (
	"os"
	"time"
)

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

func NewJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:     os.Getenv("JWT_SECRET"),
		Expiration: 24 * time.Hour,
	}
}
