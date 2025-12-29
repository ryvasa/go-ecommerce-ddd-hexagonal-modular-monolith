package config

import "os"

type PostgresConfig struct {
	DSN string
}

func LoadPostgresConfig() *PostgresConfig {
	// Railway provides DATABASE_URL automatically
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		return &PostgresConfig{DSN: databaseURL}
	}

	// Fallback to POSTGRES_DSN for local development
	return &PostgresConfig{
		DSN: GetEnv("POSTGRES_DSN", true),
	}
}
