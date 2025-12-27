package user

import (
	"github.com/google/wire"

	userHttp "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/in/http"
	userRepo "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/out/persistence"

	userService "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/service"
)

var Module = wire.NewSet(
	userRepo.NewGormUserRepository,
	userService.NewUserService,
	userHttp.NewHandler,
)
