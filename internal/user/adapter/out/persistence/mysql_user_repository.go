package persistence

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
)

type MySQLUserRepository struct{}

func (r *MySQLUserRepository) Save(ctx context.Context, user *entity.User) error {
	// implement DB logic later
	return nil
}
