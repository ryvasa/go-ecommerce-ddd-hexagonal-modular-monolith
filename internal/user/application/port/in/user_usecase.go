package in

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
)

// RegisterUserCommand represents input data required
// to register a new user in the system.
type RegisterUserCommand struct {
	Email    string
	Password string
}

// LoginUserCommand represents input data required
// to login a user in the system.
type LoginUserCommand struct {
	Email    string
	Password string
}

type UserUsecase interface {
	Register(ctx context.Context, cmd RegisterUserCommand) error
	Login(ctx context.Context, cmd LoginUserCommand) (string, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
}
