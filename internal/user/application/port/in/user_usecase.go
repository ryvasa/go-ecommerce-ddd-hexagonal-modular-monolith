package in

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
)

// RegisterUserCommand represents input data required
// to register a new user in the system.
type RegisterUserCommand struct {
	Email     string
	Password  string
	Username  string
	FirstName string
	LastName  string
}

// UserUsecase handles user management operations
type UserUsecase interface {
	Register(ctx context.Context, cmd RegisterUserCommand) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
}

// UserReader provides read-only access to user data for other modules
type UserReader interface {
	VerifyCredentials(ctx context.Context, email valueobject.Email, plainPassword string) (userID string, err error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	Exists(ctx context.Context, userID string) (bool, error)
}
