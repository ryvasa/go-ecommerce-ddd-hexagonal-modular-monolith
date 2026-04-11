---
description: full development lifecycle for implementing a new feature in go-ecommerce-ddd-hexagonal-modular-monolith
---

# HumanLayer Workflow — DDD Hexagonal Modular Monolith

Primary workflow to ensure every significant change goes through human review and
approval. Use this for tasks that span multiple layers (domain → application → adapter).

---

## 1. Research & Plan

Before writing any code:

1. Identify the affected **module** (bounded context)
2. Identify the affected **layers** (domain / application / adapter)
3. Validate against architecture rules (`.agent/rules/`)
4. Load relevant skills:
   - New entity? → Load `domain-entity` skill
   - New value object? → Load `value-object` skill
   - New endpoint? → Load `http-handler` skill
   - New persistence? → Load `gorm-persistence` skill
   - Cross-module? → Load `cross-module` skill
   - Wire changes? → Load `wire-module` skill
5. Create an `implementation_plan.md` that includes:
   - Module and layer identification
   - Files to create/modify/delete (grouped by layer)
   - Dependency direction verification
   - Cross-module boundary analysis
   - Verification plan

**[APPROVAL REQUIRED]** — Present the plan to the user and wait for approval.

---

## 2. Implement (Inside-Out Order)

After plan approval, implement in this order:

1. **Domain layer** first (entity, value object, error, repository interface)
2. **Application layer** second (port/in, port/out, service)
3. **Adapter/out layer** third (persistence model, GORM repository, security)
4. **Adapter/in layer** fourth (handler, request/response DTO, mapper, error mapping)
5. **Wire** last (module.go, wire.go)

Create a `task.md` to track implementation progress.

---

## 3. Verify

Run the `/verify` workflow for build checks and test verification.

**[APPROVAL REQUIRED]** — Report verification results to the user before committing.

---

## 4. Commit

Run the `/commit` workflow to stage, review, and commit changes.

---

## 5. Document

Create a `walkthrough.md` that documents:

- What changed and why
- Layer-by-layer summary
- Verification results
