package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
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
) *Server {
	r := gin.Default()

	r.Use(
		middleware.RequestID(),
		middleware.Logger(log),
	)

	userHandler.Register(r)
	taskHandler.Register(r)

	return &Server{Engine: r}
}
func main() {
	config.LoadEnv()

	logger := logger.NewLogger()

	// 👇 PEGANG CONCRETE
	inMemoryPublisher := eventInfra.NewInMemoryPublisher()

	// register handler (infra concern)
	taskCreatedHandler := eventhandler.NewTaskCreatedHandler(logger)
	inMemoryPublisher.Register(
		"task.created",
		func(ctx context.Context, e sharedEvent.Event) error {
			return taskCreatedHandler.Handle(ctx, e.(taskEvent.TaskCreated))
		},
	)

	// 👇 PASS SEBAGAI INTERFACE
	server, err := InitializeServer(inMemoryPublisher)
	if err != nil {
		log.Fatal(err)
	}

	server.Engine.Run(":" + config.GetEnv("APP_PORT", true))
}
