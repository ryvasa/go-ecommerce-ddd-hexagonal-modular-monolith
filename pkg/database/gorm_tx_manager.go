package database

import (
	"context"

	"gorm.io/gorm"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/transaction"
)

type GormTxManager struct {
	db *gorm.DB
}

func NewGormTxManager(db *gorm.DB) transaction.Manager {
	return &GormTxManager{db: db}
}

func (m *GormTxManager) WithinTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctxWithTx := context.WithValue(ctx, TxKey{}, tx)
		return fn(ctxWithTx)
	})
}
