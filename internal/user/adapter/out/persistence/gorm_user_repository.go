package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) out.UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Save(ctx context.Context, user *entity.User) error {
	model := UserModel{
		ID:       user.ID(),
		Email:    user.Email().Value(),
		Password: user.Password().Hash(),
	}

	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {

	var model UserModel
	if err := r.db.WithContext(ctx).
		First(&model, "email = ?", email.Value()).
		Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	emailVO, err := valueobject.NewEmail(model.Email)
	if err != nil {
		return nil, err
	}

	passwordVO, err := valueobject.NewHashedPassword(model.Password)
	if err != nil {
		return nil, err
	}

	return entity.RehydrateUser(
		model.ID,
		emailVO,
		passwordVO,
	), nil
}

func (r *GormUserRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Debug().
		Model(&UserModel{}).
		Where("id = ?", id).
		Count(&count).Error

	return count > 0, err
}

func (r *GormUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).
		First(&model, "id = ?", id).
		Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	emailVO, err := valueobject.NewEmail(model.Email)
	if err != nil {
		return nil, err
	}

	passwordVO, err := valueobject.NewHashedPassword(model.Password)
	if err != nil {
		return nil, err
	}

	return entity.RehydrateUser(
		model.ID,
		emailVO,
		passwordVO,
	), nil
}
