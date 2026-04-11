---
name: domain-error
description: >
  Creates typed domain errors with code, message, and category metadata for the
  go-ecommerce-ddd-hexagonal-modular-monolith project. Covers error struct definition,
  factory functions, HTTP status mapping, and per-module error registration. Use when
  adding new domain errors or mapping domain errors to HTTP responses.
metadata:
  author: ryvasa
  version: "1.0.0"
  domain: ddd
  triggers: domain error, error, ErrNotFound, ErrConflict, error mapping, HTTP status, error category
  role: specialist
  scope: domain-layer
  output-format: code
  related-skills: domain-entity, http-handler
---

# Domain Error — Typed Business Errors

This skill enforces domain error patterns used in the go-ecommerce-ddd-hexagonal-modular-monolith
project. Errors MUST carry business meaning with structured metadata.

## Error Architecture

```
domain/error/error.go          ← Error struct + category enum + factory functions
domain/error/{module}_error.go ← Module-specific error variables
adapter/in/http/error_mapping.go ← Domain error → HTTP status mapping
```

## Error Structure

### Base Error Type (per module)

```go
package error

// ErrorCategory represents the type of error.
type ErrorCategory string

const (
    CategoryValidation   ErrorCategory = "VALIDATION"
    CategoryNotFound     ErrorCategory = "NOT_FOUND"
    CategoryConflict     ErrorCategory = "CONFLICT"
    CategoryUnauthorized ErrorCategory = "UNAUTHORIZED"
    CategoryForbidden    ErrorCategory = "FORBIDDEN"
    CategoryBusiness     ErrorCategory = "BUSINESS_RULE"
)

// UserError represents a domain error with metadata.
type UserError struct {
    Code     string
    Message  string
    Category ErrorCategory
    Details  map[string]interface{}
}

func (e *UserError) Error() string {
    return e.Message
}

// WithDetails adds additional context to the error.
func (e *UserError) WithDetails(details map[string]interface{}) *UserError {
    e.Details = details
    return e
}

// GetCode returns the error code.
func (e *UserError) GetCode() string {
    return e.Code
}

// GetCategory returns the error category.
func (e *UserError) GetCategory() ErrorCategory {
    return e.Category
}
```

### Factory Functions

```go
func NewUserError(code, message string, category ErrorCategory) *UserError {
    return &UserError{
        Code:     code,
        Message:  message,
        Category: category,
        Details:  make(map[string]interface{}),
    }
}

func NewValidationError(code, message string) *UserError {
    return NewUserError(code, message, CategoryValidation)
}

func NewNotFoundError(code, message string) *UserError {
    return NewUserError(code, message, CategoryNotFound)
}

func NewConflictError(code, message string) *UserError {
    return NewUserError(code, message, CategoryConflict)
}

func NewUnauthorizedError(code, message string) *UserError {
    return NewUserError(code, message, CategoryUnauthorized)
}

func NewForbiddenError(code, message string) *UserError {
    return NewUserError(code, message, CategoryForbidden)
}

func NewBusinessRuleError(code, message string) *UserError {
    return NewUserError(code, message, CategoryBusiness)
}
```

### Module-Specific Error Variables

```go
// users_error.go
package error

var (
    // Authentication errors
    ErrEmailAlreadyUsed = NewConflictError(
        "USER_EMAIL_EXISTS",
        "email already used",
    )
    ErrInvalidCredential = NewUnauthorizedError(
        "USER_INVALID_CREDENTIALS",
        "invalid credentials",
    )
    ErrUserNotFound = NewNotFoundError(
        "USER_NOT_FOUND",
        "user not found",
    )

    // Validation errors
    ErrInvalidEmail = NewValidationError(
        "USER_INVALID_EMAIL",
        "invalid email format",
    )
    ErrWeakPassword = NewValidationError(
        "USER_WEAK_PASSWORD",
        "password does not meet security requirements",
    )

    // Business rule errors
    ErrUserInactive = NewForbiddenError(
        "USER_INACTIVE",
        "user account is inactive",
    )
)
```

### HTTP Error Mapping (Adapter Layer)

```go
// adapter/in/http/error_mapping.go
package http

import (
    "net/http"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/response"
    usererror "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/<module>/domain/error"
)

func init() {
    response.RegisterErrorMappings(map[error]response.ErrorMapping{
        usererror.ErrUserNotFound:     {Status: http.StatusNotFound, Code: "USER_NOT_FOUND"},
        usererror.ErrEmailAlreadyUsed: {Status: http.StatusConflict, Code: "USER_ALREADY_EXISTS"},
        usererror.ErrInvalidCredential: {Status: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS"},
        usererror.ErrInvalidEmail:     {Status: http.StatusBadRequest, Code: "VALIDATION_ERROR"},
        usererror.ErrUserInactive:     {Status: http.StatusForbidden, Code: "USER_INACTIVE"},
    })
}
```

## Error Naming Convention

- Error code: `{MODULE}_{ERROR_NAME}` (e.g., `USER_NOT_FOUND`, `CART_ITEM_NOT_FOUND`)
- Error variable: `Err{Description}` (e.g., `ErrUserNotFound`, `ErrItemNotFound`)
- Error struct: `{Module}Error` (e.g., `UserError`, `CartError`)

## Checklist (Before Adding Domain Error)

- [ ] Error variable uses `Err` prefix
- [ ] Error code uses `{MODULE}_{ERROR}` format
- [ ] Category matches HTTP semantics (NotFound, Conflict, etc.)
- [ ] Error is defined in `domain/error/{module}_error.go`
- [ ] HTTP mapping registered in `adapter/in/http/error_mapping.go`
- [ ] Error message is user-friendly (no internal details)

## Constraints

### MUST DO
- Define error variables as package-level `var`
- Use factory functions for creation
- Map ALL domain errors to HTTP status in adapter layer
- Group errors by category in the error file

### MUST NOT DO
- DO NOT use raw string errors in domain
- DO NOT map HTTP status in domain layer
- DO NOT expose internal error details to users
- DO NOT create errors without code + category
