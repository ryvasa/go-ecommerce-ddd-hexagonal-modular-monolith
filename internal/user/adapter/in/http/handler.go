package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/middleware"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/response"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/in/http/mapper"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/in/http/request"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/in"
)

type Handler struct {
	usecase   in.UserUsecase
	validator *validator.Validate
}

func NewHandler(uc in.UserUsecase, v *validator.Validate) *Handler {
	return &Handler{usecase: uc, validator: v}
}

func (h *Handler) Register(public, protected *gin.RouterGroup) {
	public = public.Group("/users")
	protected = protected.Group("/users")

	public.POST("/register", h.register)

	// contoh endpoint protected
	protected.GET("/me",
		middleware.Require("user", "get"),
		h.me,
	)
}

func (h *Handler) register(c *gin.Context) {
	var req request.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.HandleError(c, err)
		return
	}

	cmd := in.RegisterUserCommand{
		Email:     req.Email,
		Password:  req.Password,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := h.usecase.Register(c.Request.Context(), cmd); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessResponse(c, http.StatusCreated, "user registered successfully", gin.H{
		"email": cmd.Email,
	})
}

func (h *Handler) me(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	user, err := h.usecase.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": mapper.ToUserResponse(user),
	})
}
