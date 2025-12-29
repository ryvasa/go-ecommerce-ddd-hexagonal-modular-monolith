package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	cartHttp "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/adapter/in/http"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/domain/authorization"
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
	authorizer authorization.Authorizer,
	userHandler *userHttp.Handler,
	taskHandler *taskHttp.Handler,
	cartHandler *cartHttp.Handler,
	jwtVerifier middleware.JWTVerifier,
) *Server {
	r := gin.Default()

	public := r.Group("")
	protected := r.Group("")

	protected.Use(
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.JWT(jwtVerifier),
	)

	r.Use(func(c *gin.Context) {
		c.Set("authorizer", authorizer)
		c.Next()
	})

	userHandler.Register(public, protected)
	taskHandler.Register(public, protected)
	cartHandler.Register(public, protected)

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
