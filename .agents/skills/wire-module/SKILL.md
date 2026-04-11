---
name: wire-module
description: >
  Configures Wire dependency injection for modules in the go-ecommerce-ddd-hexagonal-modular-monolith
  project. Covers module.go provider sets, wire.Bind for interface mapping, cross-module
  bindings in wire.go, and wire_gen.go regeneration. Use when registering new modules,
  adding service dependencies, or wiring cross-module ports.
metadata:
  author: ryvasa
  version: "1.0.0"
  domain: di
  triggers: Wire, module, dependency injection, provider, wire.Bind, wire.go, module.go, wiring
  role: specialist
  scope: composition-root
  output-format: code
  related-skills: hexagonal-port, cross-module
---

# Wire Module — Dependency Injection Patterns

This skill enforces Wire DI patterns used in module.go and cmd/api/wire.go.

## Architecture

```
internal/<module>/module.go     ← Module provider set (Wire)
cmd/api/wire.go                 ← Composition root (all modules + cross-module bindings)
cmd/api/wire_gen.go             ← Auto-generated (never edit manually)
```

## Module Provider Set

Each module defines its Wire provider set in `module.go`:

```go
// internal/user/module.go
package user

import (
    "github.com/google/wire"

    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/in/http"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/out/persistence"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/out/security"

    sharedvalidator "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/validator"
    userReader "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/out"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/in"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/service"
)

var Module = wire.NewSet(
    // Adapters out (persistence + security)
    persistence.NewGormUserRepository,
    security.NewBcryptHasher,

    // Application service
    service.NewUserService,

    // Shared infrastructure
    sharedvalidator.NewValidator,

    // Cross-module adapter
    userReader.NewUserReaderImpl,

    // Interface bindings (port in → service)
    wire.Bind(new(in.UserUsecase), new(*service.UserService)),
    wire.Bind(new(in.UserReader), new(*service.UserService)),

    // Adapter in (HTTP handler)
    http.NewHandler,
)
```

## Composition Root (wire.go)

Cross-module bindings and infrastructure wiring:

```go
//go:build wireinject
// +build wireinject

package main

import (
    "github.com/google/wire"

    // Infrastructure
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/config"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/database"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/logger"

    // Modules
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user"

    // Cross-module bindings
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/application/port/out"
    userOut "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/out"
)

func InitializeServer(
    eventPublisher sharedEvent.Publisher,
) (*Server, error) {
    wire.Build(
        // Config
        config.LoadPostgresConfig,
        config.NewJWTConfig,

        // Infrastructure
        database.NewGormDB,
        database.NewGormTxManager,
        logger.NewLogger,

        // Authorization
        casbinInfra.ProvideEnforcer,
        casbinInfra.ProvideAuthorizer,

        // Modules
        user.Module,
        auth.Module,
        cart.Module,

        // Cross-module: cart needs user reader
        wire.Bind(
            new(out.UserReader),
            new(*userOut.UserReaderImpl),
        ),

        NewServer,
    )
    return nil, nil
}
```

## Wire Bind Rules

### Within Module (module.go)

Bind port/in interface → concrete service:

```go
wire.Bind(new(in.UserUsecase), new(*service.UserService)),
```

### Cross Module (wire.go)

Bind module A's port/out → module B's adapter/out:

```go
// Cart module needs UserReader, implemented by User module
wire.Bind(
    new(cartPortOut.UserReader),        // Consumer's interface
    new(*userAdapterOut.UserReaderImpl), // Provider's implementation
),
```

## Regenerating wire_gen.go

After modifying `wire.go` or `module.go`:

```bash
cd cmd/api && wire
```

Or:

```bash
go generate ./cmd/api/...
```

## Checklist (Before Modifying Wire)

- [ ] New provider constructor returns correct type
- [ ] `wire.Bind` maps interface → concrete implementation
- [ ] Cross-module bindings ONLY in `cmd/api/wire.go`
- [ ] Within-module bindings in `internal/<module>/module.go`
- [ ] Run `wire` to regenerate `wire_gen.go`
- [ ] Build passes after regeneration: `go build ./...`

## Constraints

### MUST DO
- Define `var Module = wire.NewSet(...)` in module.go
- Use `wire.Bind()` for interface → implementation mapping
- Cross-module binding only in composition root (wire.go)
- Regenerate `wire_gen.go` after every wire change

### MUST NOT DO
- DO NOT edit `wire_gen.go` manually
- DO NOT use Wire in domain or application layer
- DO NOT create circular module dependencies
- DO NOT bind interfaces without a clear consumer
