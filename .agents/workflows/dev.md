---
description: start local development environment for go-ecommerce-ddd-hexagonal-modular-monolith
---

# Start Local Development

Workflow to run the e-commerce backend locally with all its dependencies.

---

## 1. Check Environment

// turbo

```bash
cat .env | head -5
```

Ensure `.env` file exists with required variables.

---

## 2. Start PostgreSQL

```bash
docker compose up -d postgres
```

Wait until the database container is healthy.

---

## 3. Run Migrations

```bash
go run cmd/migrate/main.go
```

Ensure migrations complete successfully.

---

## 4. Start Application

```bash
go run cmd/api/main.go
```

Or use air for hot reload development:

```bash
air
```

---

## 5. Verify Startup

Check that the application is running:

```bash
curl http://localhost:${APP_PORT:-8080}/
```

Expected response: `{"service":"go-ecommerce-ddd-hexagonal-modular-monolith"}`
