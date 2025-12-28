package user

import (
	"github.com/google/wire"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/in/http"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/out/persistence"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/out/security"

	sharedvalidator "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/validator"
	userReader "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/in"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/service"
)

var Module = wire.NewSet(
	persistence.NewGormUserRepository,
	service.NewUserService,
	security.NewBcryptHasher,
	security.NewJWTGenerator,
	security.NewJWTVerifier,

	sharedvalidator.NewValidator,

	userReader.NewUserReaderImpl,

	wire.Bind(
		new(out.TokenGenerator),
		new(*security.JWTGenerator),
	),

	wire.Bind(new(in.UserUsecase), new(*service.UserService)),
	wire.Bind(new(in.UserReader), new(*service.UserService)),
	http.NewHandler,
)
