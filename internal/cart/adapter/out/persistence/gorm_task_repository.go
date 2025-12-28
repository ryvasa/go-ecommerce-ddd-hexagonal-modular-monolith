package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/domain/entity"
)

type GormCartRepository struct {
	db *gorm.DB
}

func NewGormCartRepository(db *gorm.DB) out.CartRepository {
	return &GormCartRepository{db: db}
}

func (r *GormCartRepository) Save(ctx context.Context, cart *entity.Cart) error {
	db := getDB(ctx, r.db)
	return db.WithContext(ctx).Create(&CartModel{
		ID:     cart.ID,
		UserID: cart.UserID,
		Title:  cart.Title,
		Done:   cart.Done,
	}).Error
}

func (r *GormCartRepository) FindAll(ctx context.Context) ([]*entity.Cart, error) {
	var models []CartModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]*entity.Cart, 0, len(models))
	for _, m := range models {
		result = append(result, &entity.Cart{
			ID:     m.ID,
			UserID: m.UserID,
			Title:  m.Title,
			Done:   m.Done,
		})
	}

	return result, nil
}
