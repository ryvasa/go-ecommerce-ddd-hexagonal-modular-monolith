package service

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth/application/port/in"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/logger"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
)

type AuthService struct {
	userAuthenticator out.UserAuthenticator
	tokenGenerator    out.TokenGenerator
	logger            logger.Logger
}

func NewAuthService(
	userAuth out.UserAuthenticator,
	tokenGen out.TokenGenerator,
	log logger.Logger,
) in.AuthUsecase {
	return &AuthService{
		userAuthenticator: userAuth,
		tokenGenerator:    tokenGen,
		logger:            log,
	}
}

func (s *AuthService) Login(ctx context.Context, cmd in.LoginCommand) (string, error) {
	s.logger.Info("authenticating user", "email", cmd.Email)

	email, err := valueobject.NewEmail(cmd.Email)
	if err != nil {
		return "", err
	}

	userID, err := s.userAuthenticator.VerifyCredentials(ctx, email, cmd.Password)
	if err != nil {
		s.logger.Error("authentication failed", "err", err)
		return "", err
	}

	// For now, hardcode "user" role. Later this can be fetched from user
	token, err := s.tokenGenerator.Generate(userID, []string{"user"})
	if err != nil {
		s.logger.Error("failed to generate token", "err", err)
		return "", err
	}

	return token, nil
}
