package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/domain/entity"
)

type GormTaskRepository struct {
	db *gorm.DB
}

func NewGormTaskRepository(db *gorm.DB) out.TaskRepository {
	return &GormTaskRepository{db: db}
}

func (r *GormTaskRepository) Save(ctx context.Context, task *entity.Task) error {
	db := getDB(ctx, r.db)
	return db.WithContext(ctx).Create(&TaskModel{
		ID:     task.ID,
		UserID: task.UserID,
		Title:  task.Title,
		Done:   task.Done,
	}).Error
}

func (r *GormTaskRepository) FindAll(ctx context.Context) ([]*entity.Task, error) {
	var models []TaskModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]*entity.Task, 0, len(models))
	for _, m := range models {
		result = append(result, &entity.Task{
			ID:    m.ID,
			Title: m.Title,
			Done:  m.Done,
		})
	}

	return result, nil
}
