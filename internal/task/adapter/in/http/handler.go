package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/response"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/application/port/in"
)

type Handler struct {
	usecase in.TaskUsecase
}

func NewHandler(uc in.TaskUsecase) *Handler {
	return &Handler{usecase: uc}
}

func (h *Handler) Register(r *gin.Engine) {
	r.POST("/tasks", h.create)
	r.GET("/tasks", h.list)
}

func (h *Handler) create(c *gin.Context) {
	var req struct {
		Title  string `json:"title"`
		UserID string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	if err := h.usecase.Create(c.Request.Context(), req.UserID, req.Title); err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func (h *Handler) list(c *gin.Context) {
	tasks, err := h.usecase.List(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}
