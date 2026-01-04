package auth

import (
	"github.com/google/wire"

	authHttp "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth/adapter/in/http"
	authOut "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth/adapter/out"
	authSecurity "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth/adapter/out/security"
	authService "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth/application/service"
)

var Module = wire.NewSet(
	authService.NewAuthService,
	authOut.NewUserAuthenticatorImpl,
	authSecurity.NewJWTGenerator, // JWT now owned by auth
	authSecurity.NewJWTVerifier,  // JWT now owned by auth
	authHttp.NewHandler,
)
