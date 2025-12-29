package config

type PostgresConfig struct {
	DSN string
}

func LoadPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		DSN: GetEnv("POSTGRES_DSN", true),
	}
}
