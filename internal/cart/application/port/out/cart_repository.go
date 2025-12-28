package out

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/domain/entity"
)

type CartRepository interface {
	Save(ctx context.Context, cart *entity.Cart) error
	FindAll(ctx context.Context) ([]*entity.Cart, error)
}
