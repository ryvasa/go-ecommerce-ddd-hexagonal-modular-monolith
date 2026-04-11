---
name: domain-entity
description: >
  Creates and modifies domain entities with DDD behavior-rich patterns for the
  go-ecommerce-ddd-hexagonal-modular-monolith project. Enforces unexported fields,
  getter methods, domain behavior methods, constructor invariants, and RehydrateEntity
  pattern. Use when creating new entities, adding entity behavior, or modifying aggregate roots.
metadata:
  author: ryvasa
  version: "1.0.0"
  domain: ddd
  triggers: entity, aggregate, domain, NewEntity, RehydrateEntity, behavior, state transition
  role: specialist
  scope: domain-layer
  output-format: code
  related-skills: value-object, domain-error
---

# Domain Entity — DDD Aggregate Root Patterns

This skill enforces domain entity patterns used in the go-ecommerce-ddd-hexagonal-modular-monolith
project. All entities MUST comply with these patterns.

## Entity Structure

Entities live in `internal/<module>/domain/entity/` and follow strict DDD rules:

```go
package entity

import (
    "github.com/google/uuid"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/<module>/domain/valueobject"
)

// User represents the User aggregate root.
type User struct {
    id       string
    email    valueobject.Email
    password valueobject.Password
    roles    map[string]struct{}
    username string
}
```

## Rules

### Field Encapsulation (NON-NEGOTIABLE)

- ALL entity fields MUST be **unexported** (lowercase)
- External access via **getter methods** only
- State changes via **domain behavior methods** only
- NEVER use public setter methods

### Constructor Pattern

Two constructors per entity:

1. **`NewEntity(...)`** — Creates a new entity with generated ID and validated invariants
2. **`RehydrateEntity(...)`** — Reconstructs from persistence (no ID generation, minimal validation)

```go
// NewUser creates a new User with validated invariants.
func NewUser(
    email valueobject.Email,
    password valueobject.Password,
    username string,
    firstName string,
    lastName string,
) *User {
    return &User{
        id:        uuid.NewString(),
        email:     email,
        username:  username,
        firstName: firstName,
        lastName:  lastName,
        password:  password,
        roles:     map[string]struct{}{"user": {}},
    }
}

// RehydrateUser reconstructs a User from persistence data.
func RehydrateUser(
    id string,
    email valueobject.Email,
    password valueobject.Password,
) *User {
    return &User{
        id:       id,
        email:    email,
        password: password,
    }
}
```

### Getter Methods

```go
func (u *User) ID() string                   { return u.id }
func (u *User) Email() valueobject.Email      { return u.email }
func (u *User) Password() valueobject.Password { return u.password }
func (u *User) Username() string              { return u.username }
```

### Domain Behavior Methods

State changes happen ONLY through behavior methods that enforce business rules:

```go
// ChangePassword validates old password and sets new one.
func (u *User) ChangePassword(
    oldPlain string,
    newPlain string,
    hasher valueobject.PasswordHasher,
) error {
    if !u.password.Verify(oldPlain, hasher) {
        return errors.New("invalid current password")
    }
    newPassword, err := valueobject.NewPasswordFromPlain(newPlain, hasher)
    if err != nil {
        return err
    }
    u.password = newPassword
    return nil
}

// AssignRole assigns a new role to the user.
func (u *User) AssignRole(role string) error {
    if _, exists := u.roles[role]; exists {
        return errors.New("role already assigned")
    }
    u.roles[role] = struct{}{}
    return nil
}
```

### ID Generation

- Use `github.com/google/uuid` for UUID-based IDs
- ID is generated in `NewEntity()`, not in persistence layer
- `RehydrateEntity()` accepts ID from persistence

## Checklist (Before Creating an Entity)

- [ ] Fields are unexported
- [ ] `NewEntity(...)` constructor validates invariants
- [ ] `RehydrateEntity(...)` constructor for persistence reconstruction
- [ ] Getter methods for all readable fields
- [ ] Behavior methods for all state transitions
- [ ] No ORM tags (`gorm:`, `json:`, `bson:`)
- [ ] No import of framework packages (gin, gorm, wire)
- [ ] Entity file lives in `internal/<module>/domain/entity/`
- [ ] Entity has corresponding unit test `entity_test.go`

## Constraints

### MUST DO
- Use value objects for validated fields (Email, Password, etc.)
- Generate UUID in `NewEntity()` constructor
- Keep entity behavior in the entity, not in service
- Receiver name matches entity: `u` for User, `o` for Order, `c` for Cart

### MUST NOT DO
- DO NOT add ORM tags to entities
- DO NOT create public setters
- DO NOT import gin, gorm, or any adapter package
- DO NOT put orchestration logic in entities (that belongs in service)
- DO NOT make entities anemic (only struct + getters)
