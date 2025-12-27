package service

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/logger"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/in"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
)

type UserService struct {
	repo   out.UserRepository
	logger logger.Logger
}

func NewUserService(repo out.UserRepository, log logger.Logger) in.UserUsecase {
	return &UserService{
		repo:   repo,
		logger: log,
	}
}

func (s *UserService) Create(ctx context.Context, email string) error {
	s.logger.Info("creating user", "email", email)
	user := &entity.User{
		ID:    "generated-id",
		Email: email,
	}

	if err := s.repo.Save(ctx, user); err != nil {
		s.logger.Error("failed to create user", "err", err)
		return err
	}

	return nil
}
