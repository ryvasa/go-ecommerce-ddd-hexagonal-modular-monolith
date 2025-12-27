package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/database"
)

func getDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(database.TxKey{}).(*gorm.DB); ok {
		return tx
	}
	return db
}
