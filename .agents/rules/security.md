---
trigger: always_on
---

# Security Rules — DDD Hexagonal Modular Monolith

## Secrets

- NEVER hardcode passwords, API keys, tokens, or secrets in source code.
- ALL secrets MUST come from environment variables (via `pkg/config`).
- NEVER commit `.env` files with real credentials.

## Input Validation

- ALWAYS validate user input in adapter/in/http handler using `validator.Struct()`.
- ALWAYS parse and validate IDs from URL params — never trust raw string input.
- ALWAYS use parameterized queries via GORM — NEVER concatenate SQL strings.
- Domain validation (business rules) happens in domain layer (entity/VO constructors).

## Authentication & Authorization

- Protected routes MUST use JWT middleware + Casbin authorization middleware.
- Public routes (login, register) go in the `public` router group.
- NEVER expose internal error messages to users (raw SQL, stack traces).
- Authorization is infrastructure concern — NEVER in domain layer.

## Password Handling

- Passwords MUST be hashed via `PasswordHasher` port (bcrypt adapter).
- Plain text passwords NEVER stored or logged.
- Password policy validated in domain (`valueobject.NewPasswordFromPlain`).

## Logging

- NEVER log sensitive data (passwords, tokens, PII).
- Log errors at the application/service level.
- Use structured logging via shared logger interface.
