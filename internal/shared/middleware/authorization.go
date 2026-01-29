package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/logger"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/authorization"
)

// Require creates a middleware that checks authorization for a specific object and action.
// It logs all authorization decisions for audit trail purposes.
func Require(obj, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizer := c.MustGet("authorizer").(authorization.Authorizer)
		log := c.MustGet("logger").(logger.Logger)
		role := c.GetString("role")
		userID := c.GetString("user_id")
		requestID := c.GetString("request_id")

		ok, err := authorizer.Enforce(role, obj, act)

		// Audit log for authorization decision
		logFields := map[string]interface{}{
			"request_id": requestID,
			"user_id":    userID,
			"role":       role,
			"object":     obj,
			"action":     act,
			"allowed":    ok,
		}

		if err != nil {
			logFields["error"] = err.Error()
			log.Error("authorization check failed", logFields)
			c.AbortWithStatusJSON(500, gin.H{
				"message": "internal server error",
			})
			return
		}

		if !ok {
			log.Info("authorization denied", logFields)
			c.AbortWithStatusJSON(403, gin.H{
				"message": "forbidden",
				"error":   "insufficient permissions",
			})
			return
		}

		log.Info("authorization granted", logFields)
		c.Next()
	}
}
