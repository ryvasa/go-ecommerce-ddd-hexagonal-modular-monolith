package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/apperror"
	sharedEvent "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/event"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/transaction"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/application/port/in"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/domain/entity"
	taskerror "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/domain/error"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/domain/event"
	userIn "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/in"
)

type TaskService struct {
	repo           out.TaskRepository
	userReader     userIn.UserReader
	txManager      transaction.Manager
	eventPublisher sharedEvent.Publisher
}

func NewTaskService(
	repo out.TaskRepository,
	userReader userIn.UserReader,
	txManager transaction.Manager,
	eventPublisher sharedEvent.Publisher,
) in.TaskUsecase {
	return &TaskService{
		repo:           repo,
		userReader:     userReader,
		txManager:      txManager,
		eventPublisher: eventPublisher,
	}
}

func (s *TaskService) Create(ctx context.Context, userID, title string) error {
	return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Validate user ID
		if strings.TrimSpace(userID) == "" {
			return taskerror.ErrEmptyUserID
		}

		// Validate title
		if strings.TrimSpace(title) == "" {
			return taskerror.ErrEmptyTitle
		}

		// Check if user exists (cross-module check, use shared apperror)
		exists, err := s.userReader.Exists(txCtx, userID)
		if err != nil {
			return err
		}
		if !exists {
			return apperror.NotFound("user not found")
		}

		task := &entity.Task{
			ID:     uuid.NewString(),
			UserID: userID,
			Title:  title,
			Done:   false,
		}

		err = s.repo.Save(txCtx, task)
		if err != nil {
			return err
		}

		err = s.eventPublisher.Publish(
			txCtx,
			event.TaskCreated{
				TaskID: task.ID,
				UserID: task.UserID,
			},
		)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println("Task created successfully")
		return nil
	})
}

func (s *TaskService) List(ctx context.Context) ([]*entity.Task, error) {
	return s.repo.FindAll(ctx)
}
