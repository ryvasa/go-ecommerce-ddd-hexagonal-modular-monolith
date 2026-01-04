package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth/adapter/in/http/request"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth/application/port/in"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/response"
)

type Handler struct {
	usecase   in.AuthUsecase
	validator *validator.Validate
}

func NewHandler(uc in.AuthUsecase, v *validator.Validate) *Handler {
	return &Handler{usecase: uc, validator: v}
}

func (h *Handler) Register(public, protected *gin.RouterGroup) {
	auth := public.Group("/auth")

	auth.POST("/login", h.login)
	// Future: /logout, /refresh, etc.
}

func (h *Handler) login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.HandleError(c, err)
		return
	}

	token, err := h.usecase.Login(c.Request.Context(), in.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessResponse(c, http.StatusOK, "login successful", gin.H{
		"access_token": token,
		"token_type":   "Bearer",
	})
}
