package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/logger"
)

func Logger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if reqID, ok := c.Get(RequestIDKey); ok {
			log = log.With("request_id", reqID)
		}

		c.Set("logger", log)
		c.Next()
	}
}
