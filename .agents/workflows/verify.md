---
description: build check, lint, and test verification for go-ecommerce-ddd-hexagonal-modular-monolith
---

# Build & Verify

Workflow to verify that code changes are valid, compile, and pass tests.

---

## 1. Build Check

// turbo

```bash
go build ./...
```

Ensure there are no compile errors before proceeding.

---

## 2. Lint Check

// turbo

```bash
go vet ./...
```

---

## 3. Run Unit Tests

// turbo

```bash
go test -race ./internal/.../domain/... ./internal/.../application/... -v -count=1
```

Run domain and application layer unit tests with race detector.

---

## 4. Run Integration Tests (if applicable)

Run adapter/out integration tests only if persistence changes were made:

```bash
go test -race ./internal/.../adapter/out/... -v -count=1
```

---

## 5. Wire Check

If `module.go` or `wire.go` was modified:

```bash
cd cmd/api && wire
```

Ensure Wire can generate `wire_gen.go` without errors.

---

## 6. Report Results

**[APPROVAL REQUIRED]** — Report verification results to the user before proceeding.

The report must include:

- Build status (pass/fail)
- Lint issues (if any)
- Test results (pass/fail per layer)
- Wire generation status (if applicable)
