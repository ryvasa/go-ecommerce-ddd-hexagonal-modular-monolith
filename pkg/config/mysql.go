package config

type MySQLConfig struct {
	DSN string
}

func LoadMySQLConfig() *MySQLConfig {
	return &MySQLConfig{
		DSN: GetEnv("MYSQL_DSN", true),
	}
}
