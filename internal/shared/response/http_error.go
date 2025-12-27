package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/apperror"
)

func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperror.Error); ok {
		status := http.StatusInternalServerError

		switch appErr.Code {
		case "NOT_FOUND":
			status = http.StatusNotFound
		case "VALIDATION_ERROR":
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{
			"error":   appErr.Code,
			"message": appErr.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   "INTERNAL_ERROR",
		"message": "unexpected error",
	})
}
