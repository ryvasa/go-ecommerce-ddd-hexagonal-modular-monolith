package out

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/domain/entity"
)

type TaskRepository interface {
	Save(ctx context.Context, task *entity.Task) error
	FindAll(ctx context.Context) ([]*entity.Task, error)
}
