package main

import (
	"log"

	"github.com/gin-gonic/gin"

	userHttp "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/in/http"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/config"
)

type Server struct {
	Engine *gin.Engine
}

func NewServer(
	userHandler *userHttp.Handler,
) *Server {
	r := gin.Default()

	userHandler.Register(r)

	return &Server{Engine: r}
}

func main() {
	config.LoadEnv()

	server, err := InitializeServer()
	if err != nil {
		log.Fatal(err)
	}

	server.Engine.Run(":" + config.GetEnv("APP_PORT", true))
}
