---
name: cross-module
description: >
  Implements cross-module communication via port/out interfaces and adapter/out
  implementations for the go-ecommerce-ddd-hexagonal-modular-monolith project.
  Use when one module needs data from another module without violating DDD boundaries.
metadata:
  author: ryvasa
  version: "1.0.0"
  domain: modular-monolith
  triggers: cross-module, module boundary, port/out, UserReader, wire.Bind
  role: specialist
  scope: cross-cutting
  output-format: code
  related-skills: hexagonal-port, wire-module
---

# Cross-Module Communication — Boundary-Safe Patterns

## The Rule

Module A NEVER imports module B's domain. Communication via port/out + adapter/out.

## Step-by-Step

### 1. Consumer defines port/out interface

```go
// internal/cart/application/port/out/user_reader.go
package out

import "context"

type UserReader interface {
    Exists(ctx context.Context, userID string) (bool, error)
}
```

### 2. Provider creates adapter/out implementation

```go
// internal/user/adapter/out/user_reader_impl.go
package out

import (
    "context"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain"
)

type UserReaderImpl struct {
    repo domain.UserRepository
}

func NewUserReaderImpl(repo domain.UserRepository) *UserReaderImpl {
    return &UserReaderImpl{repo: repo}
}

func (r *UserReaderImpl) Exists(ctx context.Context, userID string) (bool, error) {
    return r.repo.ExistsByID(ctx, userID)
}
```

### 3. Wire binding in composition root

```go
// cmd/api/wire.go
wire.Bind(
    new(cartPortOut.UserReader),
    new(*userAdapterOut.UserReaderImpl),
),
```

### 4. Consumer uses the port in service

```go
type CartService struct {
    repo       domain.CartRepository
    userReader out.UserReader
}
```

## File Organization

```
internal/cart/application/port/out/user_reader.go     <- Consumer defines interface
internal/user/adapter/out/user_reader_impl.go         <- Provider implements
internal/user/module.go                               <- Provider registers impl
cmd/api/wire.go                                       <- Cross-module binding
```

## Naming Convention

| Component | Pattern | Example |
|-----------|---------|---------|
| Consumer port | `{Entity}Reader` | `UserReader` |
| Provider impl | `{Entity}{Capability}Impl` | `UserReaderImpl` |

## Constraints

### MUST DO
- Consumer defines the contract
- Provider implements it
- Wire binding ONLY in `cmd/api/wire.go`
- Keep interfaces minimal

### MUST NOT DO
- DO NOT import `internal/<other_module>/domain/`
- DO NOT create circular module dependencies
- DO NOT put cross-module binding in module.go
