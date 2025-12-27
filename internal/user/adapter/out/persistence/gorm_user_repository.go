package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) out.UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Save(ctx context.Context, user *entity.User) error {
	model := UserModel{
		ID:    user.ID,
		Email: user.Email,
	}

	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *GormUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &entity.User{
		ID:    model.ID,
		Email: model.Email,
	}, nil
}

func (r *GormUserRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Debug().
		Model(&UserModel{}).
		Where("id = ?", id).
		Count(&count).Error

	return count > 0, err
}
