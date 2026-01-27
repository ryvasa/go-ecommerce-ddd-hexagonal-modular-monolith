package domain

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/domain/entity"
)

// CartRepository is a domain concept - collection of Cart aggregates.
// This interface belongs in domain layer, not application layer.
type CartRepository interface {
	Save(ctx context.Context, cart *entity.Cart) error
	FindAll(ctx context.Context) ([]*entity.Cart, error)
}
