package mapper

import (
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/in/http/response"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
)

func ToUserResponse(u *entity.User) response.UserResponse {
	return response.UserResponse{
		ID:     u.ID(),
		Email:  u.Email().Value(),
		Active: u.IsActive(),
		Roles:  u.Roles(),
	}
}
