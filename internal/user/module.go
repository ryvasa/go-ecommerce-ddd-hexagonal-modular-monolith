package user

import (
	"github.com/google/wire"

	userHttp "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/in/http"
	userRepo "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/out/persistence"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/in"
	userService "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/service"
)

var Module = wire.NewSet(
	userRepo.NewGormUserRepository,
	userService.NewUserService,
	// expose ports
	wire.Bind(new(in.UserUsecase), new(*userService.UserService)),
	wire.Bind(new(in.UserReader), new(*userService.UserService)),
	userHttp.NewHandler,
)
