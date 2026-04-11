---
description: add a new bounded context module to the modular monolith (domain → application → adapter → wire)
---

# New Module — Add Bounded Context

// turbo-all

Workflow to scaffold a new module (bounded context) following DDD Hexagonal patterns.

---

## 1. Create Directory Structure

```bash
MODULE=<module_name>
mkdir -p internal/$MODULE/domain/entity
mkdir -p internal/$MODULE/domain/valueobject
mkdir -p internal/$MODULE/domain/error
mkdir -p internal/$MODULE/application/port/in
mkdir -p internal/$MODULE/application/port/out
mkdir -p internal/$MODULE/application/service
mkdir -p internal/$MODULE/adapter/in/http/mapper
mkdir -p internal/$MODULE/adapter/in/http/request
mkdir -p internal/$MODULE/adapter/in/http/response
mkdir -p internal/$MODULE/adapter/out/persistence
```

Replace `<module_name>` with the actual module name (e.g., `order`, `product`, `payment`).

---

## 2. Create Domain Layer

Load the following skills and create files in order:

1. `domain-entity` skill → Create entity in `domain/entity/`
2. `value-object` skill → Create VOs in `domain/valueobject/`
3. `domain-error` skill → Create errors in `domain/error/`
4. Create `domain/repository.go` with repository interface

---

## 3. Create Application Layer

Load the following skills:

1. `hexagonal-port` skill → Create `port/in/{module}_usecase.go`
2. `hexagonal-port` skill → Create `port/out/` interfaces
3. Create `service/{module}_service.go` implementing port/in

---

## 4. Create Adapter Layer

Load the following skills:

1. `gorm-persistence` skill → Create persistence model + repository
2. `http-handler` skill → Create handler + request/response DTOs + mapper + error mapping
3. `db-migration` skill → Create SQL migration files

---

## 5. Wire Module

Load `wire-module` skill:

1. Create `internal/<module>/module.go` with Wire provider set
2. Add module to `cmd/api/wire.go`
3. Add handler to `cmd/api/main.go` NewServer function
4. Regenerate: `cd cmd/api && wire`

---

## 6. Verify

Run `/verify` workflow.

---

## 7. Commit

Run `/commit` workflow.
