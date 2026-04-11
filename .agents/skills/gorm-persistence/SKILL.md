---
name: gorm-persistence
description: >
  Implements GORM persistence adapters with explicit entity-model mapping for the
  go-ecommerce-ddd-hexagonal-modular-monolith project. Covers persistence model definition,
  GORM repository implementation, entity-to-model mapping, WithContext usage, and
  ErrRecordNotFound handling. Use when creating repositories, persistence models, or
  modifying database access patterns.
metadata:
  author: ryvasa
  version: "1.0.0"
  domain: persistence
  triggers: GORM, repository, persistence, model, database, PostgreSQL, Save, FindBy, mapping, gorm model
  role: specialist
  scope: adapter-out
  output-format: code
  related-skills: domain-entity, value-object, db-migration
---

# GORM Persistence — Repository Adapter Patterns

This skill enforces GORM persistence patterns used in the adapter/out/persistence layer.

## Architecture

```
adapter/out/persistence/
  gorm_{entity}_model.go        ← Persistence model (GORM tags)
  gorm_{entity}_repository.go   ← Repository implementation
  gorm_{entity}_repository_integration_test.go  ← Integration test
  db_helper.go                  ← Shared DB utilities
```

## Persistence Model

Models are separate from domain entities and contain GORM tags:

```go
// gorm_user_model.go
package persistence

import (
    "time"
    "gorm.io/gorm"
)

type UserModel struct {
    gorm.Model
    ID              string     `gorm:"primaryKey;size:36"`
    Email           string     `gorm:"uniqueIndex;size:255;not null"`
    EmailVerifiedAt *time.Time `gorm:"type:datetime"`
    Password        string     `gorm:"size:255;not null"`
    Username        string     `gorm:"size:255;not null"`
    FirstName       string     `gorm:"size:255;not null"`
    LastName        string     `gorm:"size:255;not null"`
    Phone           *string    `gorm:"size:255"`
    BirthDate       *time.Time `gorm:"type:datetime"`
}

func (UserModel) TableName() string {
    return "users"
}
```

## Repository Implementation

```go
// gorm_user_repository.go
package persistence

import (
    "context"
    "errors"

    "gorm.io/gorm"

    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/entity"
    "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
)

// GormUserRepository implements domain.UserRepository using GORM.
type GormUserRepository struct {
    db *gorm.DB
}

// NewGormUserRepository creates a new GORM-based UserRepository.
func NewGormUserRepository(db *gorm.DB) domain.UserRepository {
    return &GormUserRepository{db: db}
}
```

## Entity ↔ Model Mapping

### Entity → Model (Save)

```go
func (r *GormUserRepository) Save(ctx context.Context, user *entity.User) error {
    model := UserModel{
        ID:        user.ID(),
        Email:     user.Email().Value(),
        Password:  user.Password().Hash(),
        Username:  user.Username(),
        FirstName: user.FirstName(),
        LastName:  user.LastName(),
    }
    return r.db.WithContext(ctx).Create(&model).Error
}
```

### Model → Entity (Find)

```go
func (r *GormUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
    var model UserModel
    if err := r.db.WithContext(ctx).
        First(&model, "email = ?", email.Value()).
        Error; err != nil {

        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil  // Not found = nil, NOT error
        }
        return nil, err
    }

    // Reconstruct value objects
    emailVO, err := valueobject.NewEmail(model.Email)
    if err != nil {
        return nil, err
    }
    passwordVO, err := valueobject.NewHashedPassword(model.Password)
    if err != nil {
        return nil, err
    }

    return entity.RehydrateUser(
        model.ID,
        emailVO,
        passwordVO,
    ), nil
}
```

## Critical Patterns

### WithContext (NON-NEGOTIABLE)

Every GORM operation MUST use `r.db.WithContext(ctx)`:

```go
// ✅ Correct
r.db.WithContext(ctx).Create(&model).Error
r.db.WithContext(ctx).First(&model, "id = ?", id).Error
r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", id).Count(&count).Error

// ❌ Wrong
r.db.Create(&model).Error  // Missing context!
```

### ErrRecordNotFound Handling

- `FindByX` methods return `(nil, nil)` when not found — NOT an error
- Only propagate actual database errors
- NEVER return `gorm.ErrRecordNotFound` to domain

```go
if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, nil  // Let service decide what to do
}
return nil, err  // Real error
```

### Exists Pattern

```go
func (r *GormUserRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&UserModel{}).
        Where("id = ?", id).
        Count(&count).Error
    return count > 0, err
}
```

## Constructor Return Type

Constructor MUST return the **domain interface**, not the concrete struct:

```go
// ✅ Correct
func NewGormUserRepository(db *gorm.DB) domain.UserRepository {
    return &GormUserRepository{db: db}
}

// ❌ Wrong
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
    return &GormUserRepository{db: db}
}
```

## Checklist (Before Creating a Repository)

- [ ] Persistence model is separate from domain entity
- [ ] Model has GORM tags, entity does NOT
- [ ] Constructor returns domain repository interface
- [ ] All methods use `r.db.WithContext(ctx)`
- [ ] `ErrRecordNotFound` returns `nil, nil` (not error)
- [ ] Entity ↔ model mapping is explicit (no auto-mapping)
- [ ] Value objects are reconstructed properly in Find methods
- [ ] Integration test file exists

## Constraints

### MUST DO
- Always use `r.db.WithContext(ctx)` for GORM operations
- Always handle `gorm.ErrRecordNotFound` explicitly
- Always use parameterized queries (`"email = ?"`, never string concat)
- Model has explicit `TableName()` method
- Constructor returns domain interface type

### MUST NOT DO
- DO NOT return `gorm.ErrRecordNotFound` to service layer
- DO NOT use auto-mapping between entity and model
- DO NOT skip value object reconstruction in Find methods
- DO NOT import gin or HTTP packages
- DO NOT put business logic in repository
