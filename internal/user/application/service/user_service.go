package service

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/logger"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/in"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
	usererror "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/error"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
)

type UserService struct {
	repo           domain.UserRepository
	passwordHasher out.PasswordHasher
	logger         logger.Logger
}

var _ in.UserUsecase = (*UserService)(nil)
var _ in.UserReader = (*UserService)(nil)

func NewUserService(repo domain.UserRepository, hasher out.PasswordHasher, log logger.Logger) *UserService {
	return &UserService{
		repo:           repo,
		passwordHasher: hasher,
		logger:         log,
	}
}

func (s *UserService) Register(ctx context.Context, cmd in.RegisterUserCommand) error {

	s.logger.Info("registering user", "email", cmd.Email)

	email, err := valueobject.NewEmail(cmd.Email)
	if err != nil {
		return err
	}

	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existing != nil {
		return usererror.ErrEmailAlreadyUsed
	}

	password, err := valueobject.NewPasswordFromPlain(
		cmd.Password,
		s.passwordHasher,
	)
	if err != nil {
		return err
	}

	user := entity.NewUser(email, password, cmd.Username, cmd.FirstName, cmd.LastName)

	if err := s.repo.Save(ctx, user); err != nil {
		s.logger.Error("failed to register user", "err", err)
		return err
	}

	return nil
}

func (s *UserService) VerifyCredentials(ctx context.Context, email valueobject.Email, plainPassword string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", usererror.ErrInvalidCredential
	}

	if !user.Password().Verify(plainPassword, s.passwordHasher) {
		return "", usererror.ErrInvalidCredential
	}

	return user.ID(), nil
}

func (s *UserService) Exists(ctx context.Context, userID string) (bool, error) {
	return s.repo.ExistsByID(ctx, userID)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*entity.User, error) {
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, usererror.ErrUserNotFound
	}

	return data, nil
}
