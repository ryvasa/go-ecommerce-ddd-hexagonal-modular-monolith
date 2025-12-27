package eventhandler

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/logger"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/domain/event"
)

type TaskCreatedHandler struct {
	log logger.Logger
}

func NewTaskCreatedHandler(log logger.Logger) *TaskCreatedHandler {
	return &TaskCreatedHandler{log: log}
}

func (h *TaskCreatedHandler) Handle(
	ctx context.Context,
	e event.TaskCreated,
) error {
	h.log.Info(
		"task created event received",
		"task_id", e.TaskID,
		"user_id", e.UserID,
	)
	return nil
}
