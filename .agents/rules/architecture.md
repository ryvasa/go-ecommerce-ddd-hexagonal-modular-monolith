---
trigger: always_on
---

# Architecture Rules — DDD Hexagonal Modular Monolith

## Layer Dependencies (STRICT)

The dependency direction is INWARD only:

```
adapter/in (HTTP) → application (service, port) → domain (entity, VO, repo interface)
                                                         ↑
                                              adapter/out (persistence, security) implements
```

### NEVER violate these boundaries:

- adapter/in/http MUST NOT import `gorm.io/gorm` or any persistence package
- application/service MUST NOT import `github.com/gin-gonic/gin`
- application/service MUST NOT import any adapter package
- domain MUST NOT import any package outside domain (no gin, gorm, wire, external lib)
- adapter/out/persistence MUST NOT import `github.com/gin-gonic/gin`

### Cross-module boundaries:

- NEVER import `internal/<other_module>/domain/` from any module
- Cross-module call MUST go through port/out interface
- Wire binding for cross-module happens ONLY in `cmd/api/wire.go`

## Package Rules

- `internal/<module>/domain/entity/` — Business entities with behavior. No ORM tags.
- `internal/<module>/domain/valueobject/` — Immutable, validated at construction. No ORM tags.
- `internal/<module>/domain/error/` — Domain errors with category metadata.
- `internal/<module>/domain/repository.go` — Repository interface (domain concept).
- `internal/<module>/application/port/in/` — Usecase interfaces + command structs.
- `internal/<module>/application/port/out/` — Infrastructure port interfaces.
- `internal/<module>/application/service/` — Usecase implementation. Orchestration only.
- `internal/<module>/adapter/in/http/` — HTTP handlers, request/response DTOs, route registration.
- `internal/<module>/adapter/out/persistence/` — GORM repository + persistence model.
- `internal/<module>/adapter/out/security/` — Security adapters (bcrypt, JWT).
- `internal/shared/` — Cross-cutting only (error, response, middleware, logger, validator, event).

## Context Propagation

- EVERY function that does I/O MUST accept `context.Context` as the first parameter
- HTTP handlers pass `c.Request.Context()` to service methods
- Repositories use `r.db.WithContext(ctx)` for all GORM operations
- Never use `context.TODO()` in production code

## Wire (Dependency Injection)

- Wire ONLY in `cmd/api/wire.go` and `internal/<module>/module.go`
- Domain and application layers MUST NOT know about Wire
- All dependencies via constructor injection
- Use `wire.Bind()` for interface → implementation binding
