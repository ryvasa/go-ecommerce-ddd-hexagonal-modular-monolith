package service

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
)

type UserService struct {
	repo out.UserRepository
}

func NewUserService(repo out.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, email string) error {
	user := &entity.User{
		ID:    "generated-id",
		Email: email,
	}
	return s.repo.Save(ctx, user)
}
