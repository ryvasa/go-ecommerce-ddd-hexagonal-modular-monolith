package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}
}

func GetEnv(key string, required bool) string {
	value := os.Getenv(key)
	if value == "" && required {
		log.Fatalf("environment variable %s is required", key)
	}
	return value
}
