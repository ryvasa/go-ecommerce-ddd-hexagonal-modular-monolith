package out

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
}
