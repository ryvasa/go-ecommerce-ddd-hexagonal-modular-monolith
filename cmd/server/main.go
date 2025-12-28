package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	cartHttp "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/adapter/in/http"
	sharedEvent "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/event"
	_logger "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/logger"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/middleware"
	taskHttp "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/adapter/in/http"
	eventhandler "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/application/event"
	taskEvent "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/domain/event"
	userHttp "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/in/http"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/config"
	eventInfra "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/event"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/logger"
)

type Server struct {
	Engine *gin.Engine
}

func NewServer(
	log _logger.Logger,
	userHandler *userHttp.Handler,
	taskHandler *taskHttp.Handler,
	cartHandler *cartHttp.Handler,
) *Server {
	r := gin.Default()

	r.Use(
		middleware.RequestID(),
		middleware.Logger(log),
	)

	userHandler.Register(r)
	taskHandler.Register(r)
	cartHandler.Register(r)

	return &Server{Engine: r}
}
func main() {
	config.LoadEnv()

	logger := logger.NewLogger()

	inMemoryPublisher := eventInfra.NewInMemoryPublisher()

	taskCreatedHandler := eventhandler.NewTaskCreatedHandler(logger)
	inMemoryPublisher.Register(
		"task.created",
		func(ctx context.Context, e sharedEvent.Event) error {
			return taskCreatedHandler.Handle(ctx, e.(taskEvent.TaskCreated))
		},
	)

	server, err := InitializeServer(inMemoryPublisher)
	if err != nil {
		log.Fatal(err)
	}

	server.Engine.Run(":" + config.GetEnv("APP_PORT", true))
}
