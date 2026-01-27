package main

import (
	"log"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/config"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/database"
)

func main() {
	config.LoadEnv()

	postgresConfig := config.LoadPostgresConfig()
	db, err := database.NewGormDB(postgresConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("🔄 Running database migrations...")
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	log.Println("✅ Migrations completed successfully")
}
