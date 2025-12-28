package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Authorizer interface {
	Enforce(sub, obj, act string) (bool, error)
}

func Require(obj, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizer := c.MustGet("authorizer").(Authorizer)
		role := c.GetString("role")

		ok, err := authorizer.Enforce(role, obj, act)
		fmt.Println(role, obj, act)
		if err != nil || !ok {
			c.AbortWithStatusJSON(403, gin.H{
				"message": "forbidden",
			})
			return
		}

		c.Next()
	}
}
