package in

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/domain/entity"
)

type TaskUsecase interface {
	Create(ctx context.Context, userID, title string) error
	List(ctx context.Context) ([]*entity.Task, error)
}
