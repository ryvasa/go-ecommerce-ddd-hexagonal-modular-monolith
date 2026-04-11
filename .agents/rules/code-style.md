---
trigger: always_on
---

# Go Code Style Rules — DDD Hexagonal Modular Monolith

## Formatting

- Run `gofmt` on all Go files before committing.
- Use `goimports` to organize import groups: stdlib → third-party → internal.
- Maximum line length: 120 characters.

## Imports

- Group imports in this order, separated by blank lines:
  1. Standard library (`context`, `fmt`, `net/http`, etc.)
  2. Third-party (`github.com/gin-gonic/gin`, `gorm.io/gorm`, etc.)
  3. Internal packages (`github.com/ryvasa/go-ddd-hexagonal-modular-monolith/...`)

## Naming

- Use `camelCase` for unexported, `PascalCase` for exported.
- Receiver names: match the struct context (`u` for User, `s` for service, `h` for handler, `r` for repository).
- Interface names: describe behavior, NOT prefixed with `I` (use `UserRepository`, not `IUserRepository`).
- Error variables: use `err`, not `e` or `error`.
- Domain errors: prefix with `Err` (e.g. `ErrUserNotFound`, `ErrEmailAlreadyUsed`).

## Entity Style

- Entity fields MUST be unexported (private).
- Access via getter methods: `user.ID()`, `user.Email()`.
- State changes via domain methods: `user.ChangePassword(...)`, `user.AssignRole(...)`.
- Constructor validates invariants: `NewUser(...)`, `RehydrateUser(...)`.

## Value Object Style

- Wrap primitives in structs with unexported `value` field.
- Constructor validates and returns `(VO, error)`: `NewEmail(v string) (Email, error)`.
- Access via `.Value()` method.
- NEVER expose primitive directly.

## Comments

- Write comments in English.
- All exported functions, types, and interfaces MUST have doc comments.
- Comment format: `// FunctionName does X.` (start with the identifier name).

## Error Handling

- NEVER use `_` to discard errors without a comment justifying why.
- NEVER use `panic()` for normal error handling.
- ALWAYS check and handle errors immediately after the call that returns them.
- Domain errors MUST use typed errors with `Code` + `Category` metadata.
