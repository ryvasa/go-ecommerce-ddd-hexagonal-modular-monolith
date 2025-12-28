package service

import (
	"context"
	"fmt"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/logger"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/in"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
	usererror "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/error"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
)

type UserService struct {
	repo           out.UserRepository
	passwordHasher out.PasswordHasher
	logger         logger.Logger
	tokenGenerator out.TokenGenerator
}

var _ in.UserUsecase = (*UserService)(nil)
var _ in.UserReader = (*UserService)(nil)

func NewUserService(repo out.UserRepository, hasher out.PasswordHasher, tokenGen out.TokenGenerator, log logger.Logger) *UserService {
	return &UserService{
		repo:           repo,
		passwordHasher: hasher,
		tokenGenerator: tokenGen,
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

	user := entity.NewUser(email, password)

	if err := s.repo.Save(ctx, user); err != nil {
		s.logger.Error("failed to register user", "err", err)
		return err
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, cmd in.LoginUserCommand) (string, error) {

	s.logger.Info("user login attempt", "email", cmd.Email)

	email, err := valueobject.NewEmail(cmd.Email)
	if err != nil {
		return "", usererror.ErrInvalidCredential
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", usererror.ErrInvalidCredential
	}

	if !user.Password().Verify(cmd.Password, s.passwordHasher) {
		return "", usererror.ErrInvalidCredential
	}

	token, err := s.tokenGenerator.Generate(user.ID(), "user")
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) Exists(ctx context.Context, userID string) (bool, error) {
	return s.repo.ExistsByID(ctx, userID)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*entity.User, error) {
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	fmt.Println(data)
	return data, nil
}
