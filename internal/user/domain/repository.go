package domain

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
)

// UserRepository is a domain concept - collection of User aggregates.
// This interface belongs in domain layer, not application layer.
type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
}
