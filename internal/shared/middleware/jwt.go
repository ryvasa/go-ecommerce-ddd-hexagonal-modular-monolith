package middleware

import "github.com/gin-gonic/gin"

type JWTVerifier interface {
	Verify(token string) (Claims, error)
}

type Claims struct {
	UserID string
	Role   string
}

func JWT(verifier JWTVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatus(401)
			return
		}

		claims, err := verifier.Verify(auth)
		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
