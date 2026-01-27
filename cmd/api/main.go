package main

import (
	"log"

	"github.com/gin-gonic/gin"
	authHttp "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth/adapter/in/http"
	cartHttp "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/adapter/in/http"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/domain/authorization"
	sharedLogger "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/logger"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/middleware"
	userHttp "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/in/http"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/config"
	eventInfra "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/event"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/logger"
)

type Server struct {
	Engine *gin.Engine
}

func NewServer(
	log sharedLogger.Logger,
	authorizer authorization.Authorizer,
	userHandler *userHttp.Handler,
	authHandler *authHttp.Handler,
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
	cartHandler.Register(public, protected)

	return &Server{Engine: r}
}

func main() {
	config.LoadEnv()

	_ = logger.NewLogger()

	inMemoryPublisher := eventInfra.NewInMemoryPublisher()

	server, err := InitializeServer(inMemoryPublisher)
	if err != nil {
		log.Fatal(err)
	}

	server.Engine.Run(":" + config.GetEnv("APP_PORT", true))
}
