package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/domain/authorization"
)

func Require(obj, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizer := c.MustGet("authorizer").(authorization.Authorizer)
		role := c.GetString("role")

		ok, err := authorizer.Enforce(role, obj, act)
		if err != nil || !ok {
			c.AbortWithStatusJSON(403, gin.H{
				"message": "forbidden",
			})
			return
		}

		c.Next()
	}
}
