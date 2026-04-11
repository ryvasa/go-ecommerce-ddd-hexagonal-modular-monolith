---
name: value-object
description: >
  Creates immutable value objects with validation at construction for the
  go-ecommerce-ddd-hexagonal-modular-monolith project. Covers wrapped primitives,
  comparison semantics, hasher integration, and policy validation. Use when adding
  new validated fields like Email, Password, Money, Address, or Phone.
metadata:
  author: ryvasa
  version: "1.0.0"
  domain: ddd
  triggers: value object, VO, Email, Password, Money, Address, Phone, immutable, validation
  role: specialist
  scope: domain-layer
  output-format: code
  related-skills: domain-entity, domain-error
---

# Value Object — Immutable Validated Primitives

This skill enforces value object patterns used in the go-ecommerce-ddd-hexagonal-modular-monolith
project. All VOs MUST be immutable and validated at construction.

## Value Object Structure

VOs live in `internal/<module>/domain/valueobject/` and wrap primitives:

### Simple Value Object (e.g., Email)

```go
package valueobject

import (
    "errors"
    "regexp"
    "strings"
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// Email represents a validated email address.
type Email struct {
    value string
}

// NewEmail creates a validated Email value object.
func NewEmail(v string) (Email, error) {
    v = strings.TrimSpace(strings.ToLower(v))
    if !emailRegex.MatchString(v) {
        return Email{}, errors.New("invalid email format")
    }
    return Email{value: v}, nil
}

// Value returns the underlying email string.
func (e Email) Value() string {
    return e.value
}
```

### Complex Value Object with Infrastructure Port (e.g., Password)

When a VO needs infrastructure (like hashing), define a port interface in the same package:

```go
// PasswordHasher is a domain port for password hashing.
type PasswordHasher interface {
    Hash(plain string) (string, error)
    Compare(hash, plain string) bool
}
```

```go
package valueobject

// Password represents a hashed password.
type Password struct {
    hash string
}

// NewPasswordFromPlain creates a Password from plain text with policy validation.
func NewPasswordFromPlain(
    plain string,
    hasher PasswordHasher,
) (Password, error) {
    if err := validatePasswordPolicy(plain); err != nil {
        return Password{}, err
    }
    hash, err := hasher.Hash(plain)
    if err != nil {
        return Password{}, err
    }
    return Password{hash: hash}, nil
}

// NewHashedPassword reconstructs a Password from an existing hash.
func NewHashedPassword(hash string) (Password, error) {
    if hash == "" {
        return Password{}, errors.New("password hash cannot be empty")
    }
    return Password{hash: hash}, nil
}

// Hash returns the password hash.
func (p Password) Hash() string {
    return p.hash
}

// Verify checks if a plain text matches the hashed password.
func (p Password) Verify(plain string, hasher PasswordHasher) bool {
    return hasher.Compare(p.hash, plain)
}
```

### Policy Validation (Separate File)

Keep business rules in their own file:

```go
// password_policy.go
package valueobject

import "errors"

func validatePasswordPolicy(p string) error {
    if len(p) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    return nil
}
```

## File Organization

```
internal/<module>/domain/valueobject/
  email.go              ← Email VO
  email_test.go         ← Email unit test
  password.go           ← Password VO
  password_test.go      ← Password unit test
  password_hasher.go    ← PasswordHasher interface (domain port)
  password_policy.go    ← Password business rules
```

## Rules

### Immutability (NON-NEGOTIABLE)
- VO fields MUST be **unexported**
- NO methods that mutate state
- To "change" a VO, create a new instance

### Constructor Pattern
- Constructor MUST return `(VO, error)` — validation happens here
- Two constructors for persistence VOs:
  1. `NewXFromPlain(...)` — Creates from raw input with full validation
  2. `NewXFromStorage(...)` or `NewHashedX(...)` — Reconstructs from stored value

### Access Pattern
- Single access method: `.Value()` for simple VOs
- Named accessors for complex VOs: `.Hash()`, `.Currency()`

## Checklist (Before Creating a VO)

- [ ] Struct has unexported field(s)
- [ ] Constructor validates and returns `(VO, error)`
- [ ] Access via `.Value()` or named getter
- [ ] No mutation methods
- [ ] No ORM tags
- [ ] Unit test covers valid + invalid cases
- [ ] If needs infrastructure, define port interface in same package

## Constraints

### MUST DO
- Validate in constructor — invalid state is impossible post-construction
- Return `error` from constructor — never panic
- Use receiver name matching VO: `e` for Email, `p` for Password
- Precompile regex as package-level `var`

### MUST NOT DO
- DO NOT expose primitive fields directly
- DO NOT create mutation methods
- DO NOT import framework packages
- DO NOT skip validation for "trusted" input
