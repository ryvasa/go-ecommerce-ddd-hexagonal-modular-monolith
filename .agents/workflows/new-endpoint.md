---
description: add a new REST API endpoint to an existing module (handler → service → domain)
---

# New Endpoint — Add API Endpoint

Workflow to add a new REST API endpoint to an existing module.

---

## 1. Identify Scope

Determine:
- Which **module** owns this endpoint?
- Is it **public** or **protected** (JWT + Casbin)?
- What HTTP method and path? (e.g., `POST /api/orders`)
- Does it require new domain logic or reuse existing?

---

## 2. Domain Layer (if new logic needed)

1. Add entity behavior methods if needed → `domain-entity` skill
2. Add new value objects if needed → `value-object` skill
3. Add new domain errors if needed → `domain-error` skill
4. Update repository interface if new query needed

---

## 3. Application Layer

1. Add command/query struct to `port/in/{module}_usecase.go`
2. Add method to usecase interface
3. Implement in `service/{module}_service.go`
4. Add port/out interface if new infrastructure needed

---

## 4. Adapter Layer

1. Create request DTO in `adapter/in/http/request/`
2. Create response DTO in `adapter/in/http/response/` (if needed)
3. Add mapper in `adapter/in/http/mapper/` (if needed)
4. Add handler method in `adapter/in/http/handler.go`
5. Register route in `Register()` method
6. Register error mappings in `error_mapping.go` (if new errors)
7. Add persistence methods if needed → `gorm-persistence` skill
8. Add migration if schema change needed → `db-migration` skill

---

## 5. Wire (if new dependencies)

If new providers were added:
1. Update `module.go`
2. Update `wire.go` if cross-module
3. Regenerate: `cd cmd/api && wire`

---

## 6. Verify

Run `/verify` workflow.

---

## 7. Commit

Run `/commit` workflow.
