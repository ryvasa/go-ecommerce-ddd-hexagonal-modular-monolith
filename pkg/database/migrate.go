package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// RunMigrations runs SQL migrations from migrations/ directory using golang-migrate.
// This is the recommended approach for production.
func RunMigrations(gormDB *gorm.DB) error {
	// Get underlying sql.DB from gorm
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Create postgres driver instance
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Path to migration files
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run all pending migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("✅ No new migrations to apply")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("✅ Migrations completed successfully")
	return nil
}

// RollbackMigration rolls back the last migration.
// Use with caution!
func RollbackMigration(gormDB *gorm.DB, steps int) error {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Rollback specified number of steps
	if err := m.Steps(-steps); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	fmt.Printf("✅ Rolled back %d migration(s)\n", steps)
	return nil
}
