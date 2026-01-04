package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	authHttp "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth/adapter/in/http"
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
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/database"
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
	authHandler *authHttp.Handler,
	taskHandler *taskHttp.Handler,
	cartHandler *cartHttp.Handler,
	jwtVerifier middleware.JWTVerifier,
) *Server {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("authorizer", authorizer)
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"service": "go-ecommerce-ddd-hexagonal-modular-monolith"})
	})

	api := r.Group("/api")

	public := api.Group("")
	protected := api.Group("")

	protected.Use(
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.JWT(jwtVerifier),
	)

	userHandler.Register(public, protected)
	authHandler.Register(public, protected)
	taskHandler.Register(public, protected)
	cartHandler.Register(public, protected)

	return &Server{Engine: r}
}

func main() {
	config.LoadEnv()

	logger := logger.NewLogger()

	// Initialize database connection early for migrations
	postgresConfig := config.LoadPostgresConfig()
	db, err := database.NewGormDB(postgresConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations
	log.Println("🔄 Running database migrations...")
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	log.Println("✅ Migrations completed successfully")

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
