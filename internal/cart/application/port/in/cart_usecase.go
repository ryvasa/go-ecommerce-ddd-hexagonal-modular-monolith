package in

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/domain/entity"
)

type CartUsecase interface {
	Create(ctx context.Context, userID, title string) error
	List(ctx context.Context) ([]*entity.Cart, error)
}
