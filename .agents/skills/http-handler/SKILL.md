---
name: http-handler
description: >
  Creates Gin HTTP handlers with request validation, response formatting, error mapping,
  and route registration for the go-ecommerce-ddd-hexagonal-modular-monolith project.
  Covers handler struct, request/response DTOs, entity-to-response mapper, error mapping
  via init(), and public/protected route registration. Use when adding API endpoints
  or modifying HTTP delivery layer.
metadata:
  author: ryvasa
  version: "1.0.0"
  domain: http
  triggers: handler, endpoint, API, route, request, response, HTTP, REST, Gin, middleware
  role: specialist
  scope: adapter-in
  output-format: code
  related-skills: hexagonal-port, domain-error
---

# HTTP Handler — Gin Delivery Adapter Patterns

This skill enforces HTTP handler patterns used in the adapter/in/http layer.

## Handler Architecture

```
adapter/in/http/
  handler.go           ← Handler struct + route registration + handler methods
  error_mapping.go     ← Domain error → HTTP status mapping (init)
  mapper/              ← Entity → Response DTO mappers
    user_mapper.go
  request/             ← Request DTOs with validation tags
    user_request.go
  response/            ← Response DTOs with json tags
    user_response.go
```

## Handler Struct

```go
package http

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/middleware"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/response"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/<module>/adapter/in/http/mapper"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/<module>/adapter/in/http/request"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/<module>/application/port/in"
)

// Handler handles HTTP requests for the user module.
type Handler struct {
    usecase   in.UserUsecase
    validator *validator.Validate
}

// NewHandler creates a new HTTP handler.
func NewHandler(uc in.UserUsecase, v *validator.Validate) *Handler {
    return &Handler{usecase: uc, validator: v}
}
```

## Route Registration Pattern

```go
// Register registers public and protected routes.
func (h *Handler) Register(public, protected *gin.RouterGroup) {
    public = public.Group("/users")
    protected = protected.Group("/users")

    // Public endpoints (no auth required)
    public.POST("/register", h.register)

    // Protected endpoints (JWT + Casbin required)
    protected.GET("/me",
        middleware.Require("user", "get"),
        h.me,
    )
}
```

## Handler Method Pattern

```go
func (h *Handler) register(c *gin.Context) {
    // 1. Parse request
    var req request.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.HandleError(c, err)
        return
    }

    // 2. Validate
    if err := h.validator.Struct(req); err != nil {
        response.HandleError(c, err)
        return
    }

    // 3. Map to command
    cmd := in.RegisterUserCommand{
        Email:     req.Email,
        Password:  req.Password,
        Username:  req.Username,
        FirstName: req.FirstName,
        LastName:  req.LastName,
    }

    // 4. Call usecase with request context
    if err := h.usecase.Register(c.Request.Context(), cmd); err != nil {
        response.HandleError(c, err)
        return
    }

    // 5. Return response
    response.SuccessResponse(c, http.StatusCreated, "user registered successfully", gin.H{
        "email": cmd.Email,
    })
}
```

## Request DTO

```go
// request/user_request.go
package request

// CreateUserRequest represents the HTTP request body for user registration.
type CreateUserRequest struct {
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8"`
    Username  string `json:"username" validate:"required,min=3"`
    FirstName string `json:"first_name" validate:"required"`
    LastName  string `json:"last_name" validate:"required"`
}
```

## Response DTO & Mapper

```go
// response/user_response.go
package response

// UserResponse represents the HTTP response for user data.
type UserResponse struct {
    ID       string `json:"id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    Name     string `json:"name"`
}
```

```go
// mapper/user_mapper.go
package mapper

import (
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/adapter/in/http/response"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
)

// ToUserResponse maps a User entity to UserResponse DTO.
func ToUserResponse(user *entity.User) response.UserResponse {
    return response.UserResponse{
        ID:       user.ID(),
        Email:    user.Email().Value(),
        Username: user.Username(),
        Name:     user.FirstName() + " " + user.LastName(),
    }
}
```

## Error Mapping (init pattern)

```go
// error_mapping.go
package http

import (
    "net/http"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/response"
    usererror "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/error"
)

func init() {
    response.RegisterErrorMappings(map[error]response.ErrorMapping{
        usererror.ErrUserNotFound:     {Status: http.StatusNotFound, Code: "USER_NOT_FOUND"},
        usererror.ErrEmailAlreadyUsed: {Status: http.StatusConflict, Code: "USER_ALREADY_EXISTS"},
        usererror.ErrInvalidCredential: {Status: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS"},
        usererror.ErrInvalidEmail:     {Status: http.StatusBadRequest, Code: "VALIDATION_ERROR"},
    })
}
```

## Response Helpers

Use shared response package:

```go
// Success
response.SuccessResponse(c, http.StatusOK, "success message", data)
response.SuccessResponse(c, http.StatusCreated, "created", data)

// Error
response.HandleError(c, err)  // Auto-maps domain error → HTTP status
```

## Getting Auth Context

```go
func (h *Handler) me(c *gin.Context) {
    userID := c.MustGet("user_id").(string)  // Set by JWT middleware
    // ...
}
```

## Checklist (Before Creating a Handler)

- [ ] Handler depends on port/in usecase interface, not service directly
- [ ] Request DTOs in `request/` package with validation tags
- [ ] Response DTOs in `response/` package with json tags
- [ ] Entity → Response mapper in `mapper/` package
- [ ] Error mapping registered via `init()` in `error_mapping.go`
- [ ] Routes registered in `Register(public, protected)` method
- [ ] Handler passes `c.Request.Context()` to usecase
- [ ] Uses `response.SuccessResponse()` and `response.HandleError()`

## Constraints

### MUST DO
- Pass `c.Request.Context()` to usecase methods
- Use `response.HandleError(c, err)` for all errors
- Use `response.SuccessResponse(c, ...)` for all successes
- Register error mappings via `init()` in adapter layer
- Use `middleware.Require("resource", "action")` for Casbin

### MUST NOT DO
- DO NOT put business logic in handler
- DO NOT import gorm or persistence packages
- DO NOT call repository directly
- DO NOT return raw error messages to client
- DO NOT skip request validation
