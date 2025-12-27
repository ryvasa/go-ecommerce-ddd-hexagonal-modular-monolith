package task

import (
	"github.com/google/wire"

	taskHttp "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/adapter/in/http"
	taskRepo "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/adapter/out/persistence"

	taskService "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/task/application/service"
)

var Module = wire.NewSet(
	taskRepo.NewGormTaskRepository,
	taskService.NewTaskService,
	taskHttp.NewHandler,
)
