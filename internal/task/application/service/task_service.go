package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/transaction"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/application/port/in"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/domain/entity"
	userIn "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/in"
)

type TaskService struct {
	repo       out.TaskRepository
	userReader userIn.UserReader
	txManager  transaction.Manager
}

func NewTaskService(
	repo out.TaskRepository,
	userReader userIn.UserReader,
	txManager transaction.Manager,
) in.TaskUsecase {
	return &TaskService{
		repo:       repo,
		userReader: userReader,
		txManager:  txManager,
	}
}

func (s *TaskService) Create(ctx context.Context, userID, title string) error {
	return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		exists, err := s.userReader.Exists(txCtx, userID)
		if err != nil {
			return err
		}
		if !exists {
			return errors.New("user not found")
		}

		task := &entity.Task{
			ID:     uuid.NewString(),
			UserID: userID,
			Title:  title,
			Done:   false,
		}
		return s.repo.Save(txCtx, task)
	})
}

func (s *TaskService) List(ctx context.Context) ([]*entity.Task, error) {
	return s.repo.FindAll(ctx)
}
