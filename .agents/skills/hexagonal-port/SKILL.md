---
name: hexagonal-port
description: >
  Defines port interfaces (in and out) for the hexagonal architecture in the
  go-ecommerce-ddd-hexagonal-modular-monolith project. Covers usecase interfaces,
  command/query DTOs, infrastructure port interfaces, and repository interfaces.
  Use when defining new usecase contracts, adding infrastructure ports, or creating
  cross-module communication boundaries.
metadata:
  author: ryvasa
  version: "1.0.0"
  domain: hexagonal
  triggers: port, usecase, port in, port out, interface, command, query, repository interface
  role: specialist
  scope: application-layer
  output-format: code
  related-skills: domain-entity, wire-module, cross-module
---

# Hexagonal Port — Interface Design Patterns

This skill enforces port interface patterns for the hexagonal architecture.

## Port Architecture Overview

```
application/port/in/   ← Usecase interfaces (what the module offers)
application/port/out/  ← Infrastructure interfaces (what the module needs)
domain/repository.go   ← Repository interface (domain concept)
```

## Port In — Usecase Interfaces

Port In defines what the module offers to the outside world:

```go
// application/port/in/user_usecase.go
package in

import (
    "context"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
)

// RegisterUserCommand represents input data required
// to register a new user in the system.
type RegisterUserCommand struct {
    Email     string
    Password  string
    Username  string
    FirstName string
    LastName  string
}

// UserUsecase handles user management operations.
type UserUsecase interface {
    Register(ctx context.Context, cmd RegisterUserCommand) error
    GetByID(ctx context.Context, id string) (*entity.User, error)
}

// UserReader provides read-only access to user data for other modules.
type UserReader interface {
    VerifyCredentials(ctx context.Context, email valueobject.Email, plainPassword string) (userID string, err error)
    GetByID(ctx context.Context, id string) (*entity.User, error)
    Exists(ctx context.Context, userID string) (bool, error)
}
```

### Command vs Query Separation

- **Command structs** for write operations: `RegisterUserCommand`, `AddToCartCommand`
- **Query/Reader interfaces** for read operations: `UserReader`, `CartReader`
- Name commands as `{Action}{Entity}Command`
- Name query interfaces as `{Entity}Reader`

## Port Out — Infrastructure Interfaces

Port Out defines what the module needs from infrastructure:

```go
// application/port/out/password_hasher.go
package out

// PasswordHasher provides password hashing capabilities.
type PasswordHasher interface {
    Hash(plain string) (string, error)
    Compare(hash, plain string) bool
}
```

```go
// application/port/out/token_generator.go
package out

// TokenGenerator provides JWT token generation.
type TokenGenerator interface {
    Generate(userID string, roles []string) (string, error)
}
```

```go
// application/port/out/user_repository.go
package out

import (
    "context"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
)

// UserRepository defines persistence operations for User.
type UserRepository interface {
    Save(ctx context.Context, user *entity.User) error
    FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error)
    ExistsByID(ctx context.Context, id string) (bool, error)
    GetByID(ctx context.Context, id string) (*entity.User, error)
}
```

## Domain Repository Interface

Repository interface is a **domain concept** and lives in `domain/repository.go`:

```go
// domain/repository.go
package domain

import (
    "context"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
)

// UserRepository is a domain concept - collection of User aggregates.
type UserRepository interface {
    Save(ctx context.Context, user *entity.User) error
    FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error)
    ExistsByID(ctx context.Context, id string) (bool, error)
    GetByID(ctx context.Context, id string) (*entity.User, error)
}
```

## Compile-Time Interface Check

Service MUST verify it implements the port:

```go
var _ in.UserUsecase = (*UserService)(nil)
var _ in.UserReader = (*UserService)(nil)
```

## File Naming

| Type | File | Package |
|------|------|---------|
| Usecase port | `{module}_usecase.go` | `port/in` |
| Reader port | `{module}_reader.go` or in same file | `port/in` |
| Repository port | `user_repository.go` | `port/out` |
| Infrastructure port | `{capability}.go` | `port/out` |
| Domain repository | `repository.go` | `domain` |

## Checklist (Before Creating a Port)

- [ ] Port in: Interface describes behavior, not implementation
- [ ] Port in: Command structs for write operations
- [ ] Port out: Interface defined by consumer, not implementor
- [ ] Every method accepts `context.Context` as first parameter
- [ ] No framework imports in port interfaces
- [ ] Repository works with domain entities, not DTOs

## Constraints

### MUST DO
- Define port in same package as consumer
- Use `context.Context` as first parameter for all I/O methods
- Keep interfaces small (1-5 methods)
- Name interfaces by behavior: `UserReader`, `PasswordHasher`, `TokenGenerator`

### MUST NOT DO
- DO NOT define port in adapter package
- DO NOT return DTOs from repository interface
- DO NOT import adapter packages in port definitions
- DO NOT create fat interfaces (split by responsibility)
