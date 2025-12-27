package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/config"
)

func NewGormDB(cfg *config.MySQLConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
