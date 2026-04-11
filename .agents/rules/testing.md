---
trigger: manual
---

# Testing Rules — DDD Hexagonal Modular Monolith

---

## Testing Strategy (Layer-Based)

### 1. Domain Layer — Pure Unit Test

- Test entity behavior (state transitions, invariant enforcement)
- Test value object validation (constructor accepts/rejects)
- Test domain error creation and categorization
- NO mocks, NO database, NO framework
- Domain MUST be testable with zero external dependencies

### 2. Application Layer — Mock Port Unit Test

- Test usecase orchestration flow
- Mock ALL port/out interfaces (repository, hasher, token generator)
- Verify correct domain methods are called
- Verify error propagation from domain
- Use `testify/mock` or `gomock`

### 3. Adapter/Out Layer — Integration Test

- Test GORM repository with real database (testcontainers or test DB)
- Test entity ↔ persistence model mapping correctness
- Test error translation (gorm.ErrRecordNotFound → nil, not domain error)
- Test security adapters (bcrypt hash/compare, JWT generate/verify)

### 4. Adapter/In Layer — HTTP Handler Test

- Test request parsing and validation
- Test response format (status code, body structure)
- Test error mapping (domain error → HTTP status)
- Mock usecase interface

---

## Test File Placement

```
internal/<module>/domain/entity/
  user.go
  user_test.go             ← Domain unit test

internal/<module>/domain/valueobject/
  email.go
  email_test.go            ← Value object unit test

internal/<module>/application/service/
  user_service.go
  user_service_test.go     ← Application unit test (mock ports)

internal/<module>/adapter/out/persistence/
  gorm_user_repository.go
  gorm_user_repository_integration_test.go  ← Integration test

internal/<module>/adapter/out/security/
  bcrypt_hasher.go
  bcrypt_hasher_test.go    ← Security adapter test
```

---

## Naming Convention

- `Test<Function>_<Scenario>`
- Examples:
  - `TestNewEmail_ValidEmail`
  - `TestNewEmail_InvalidFormat`
  - `TestRegister_EmailAlreadyUsed`
  - `TestRegister_Success`
  - `TestSave_DuplicateEmail`

---

## Test Structure (AAA Pattern)

```go
func TestRegister_Success(t *testing.T) {
    // Arrange
    repo := new(MockUserRepository)
    hasher := new(MockPasswordHasher)
    svc := service.NewUserService(repo, hasher, logger)

    repo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, nil)
    repo.On("Save", mock.Anything, mock.Anything).Return(nil)
    hasher.On("Hash", mock.Anything).Return("hashed", nil)

    // Act
    err := svc.Register(ctx, cmd)

    // Assert
    require.NoError(t, err)
    repo.AssertCalled(t, "Save", mock.Anything, mock.Anything)
}
```

---

## Compile-Time Interface Check

Every service MUST have compile-time interface checks:

```go
var _ in.UserUsecase = (*UserService)(nil)
var _ in.UserReader = (*UserService)(nil)
```

---

## Anti-Patterns (FORBIDDEN)

- Testing implementation details instead of behavior
- Using real DB in domain or application unit tests
- Mocking domain entities
- Global shared mocks directory
- Ignoring failure scenarios
- Flaky tests dependent on timing
- No cleanup after integration tests
- Skipping `-race` flag

---

## CI Requirements

- All tests must pass before merge
- Run with race detector: `go test -race ./...`
- Coverage must include critical domain paths
