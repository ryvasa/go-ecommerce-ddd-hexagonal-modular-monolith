package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/logger"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/in"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
)

type UserService struct {
	repo   out.UserRepository
	logger logger.Logger
}

var (
	_ in.UserUsecase = (*UserService)(nil)
	_ in.UserReader  = (*UserService)(nil)
)

func NewUserService(repo out.UserRepository, log logger.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: log,
	}
}

func (s *UserService) Create(ctx context.Context, email string) error {
	s.logger.Info("creating user", "email", email)
	user := &entity.User{
		ID:    uuid.NewString(),
		Email: email,
	}

	if err := s.repo.Save(ctx, user); err != nil {
		s.logger.Error("failed to create user", "err", err)
		return err
	}

	return nil
}

func (s *UserService) Exists(ctx context.Context, userID string) (bool, error) {
	return s.repo.ExistsByID(ctx, userID)
}
