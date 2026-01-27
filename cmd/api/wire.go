//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	sharedEvent "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/event"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/config"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/database"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/logger"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user"

	authSecurity "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth/adapter/out/security"
	authPortOut "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/application/port/out"
	userOut "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/out"
	casbinInfra "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/authorization/casbin"
)

func InitializeServer(
	eventPublisher sharedEvent.Publisher,
) (*Server, error) {
	wire.Build(
		// config
		config.LoadPostgresConfig,
		config.NewJWTConfig,

		// infra
		database.NewGormDB,
		database.NewGormTxManager,
		logger.NewLogger,

		// authorization
		casbinInfra.ProvideEnforcer,
		casbinInfra.ProvideAuthorizer,

		// modules
		user.Module,
		auth.Module,
		cart.Module,

		// Bind UserReader for other module
		wire.Bind(
			new(out.UserReader),
			new(*userOut.UserReaderImpl),
		),

		// Bind TokenGenerator for auth module
		wire.Bind(
			new(authPortOut.TokenGenerator),
			new(*authSecurity.JWTGenerator), // Now from auth module
		),

		NewServer,
	)
	return nil, nil
}
