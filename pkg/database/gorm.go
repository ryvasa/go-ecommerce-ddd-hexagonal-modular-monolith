package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/config"
)

func NewGormDB(cfg *config.PostgresConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
