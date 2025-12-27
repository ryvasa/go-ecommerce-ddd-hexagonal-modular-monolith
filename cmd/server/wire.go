//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/config"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/database"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/logger"
)

func InitializeServer() (*Server, error) {
	wire.Build(
		// infra & config
		config.LoadMySQLConfig,
		database.NewGormDB,
		database.NewGormTxManager,
		logger.NewLogger,

		// modules
		user.Module,
		task.Module,

		// server
		NewServer,
	)
	return nil, nil
}
