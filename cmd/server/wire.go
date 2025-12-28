//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	sharedEvent "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/event"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/config"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/database"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/logger"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/application/port/out"
	userOut "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/out"
)

func InitializeServer(
	eventPublisher sharedEvent.Publisher,
) (*Server, error) {
	wire.Build(
		config.LoadMySQLConfig,
		config.NewJWTConfig,

		// infra
		database.NewGormDB,
		database.NewGormTxManager,
		logger.NewLogger,

		// modules
		user.Module,
		task.Module,
		cart.Module,

		wire.Bind(
			new(out.UserReader),
			new(*userOut.UserReaderImpl),
		),

		NewServer,
	)
	return nil, nil
}
