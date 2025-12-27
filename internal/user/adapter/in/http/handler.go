package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/in"
)

type Handler struct {
	usecase in.UserUsecase
}

func NewHandler(uc in.UserUsecase) *Handler {
	return &Handler{usecase: uc}
}

func (h *Handler) Register(r *gin.Engine) {
	r.POST("/users", h.createUser)
}

func (h *Handler) createUser(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.Create(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}
