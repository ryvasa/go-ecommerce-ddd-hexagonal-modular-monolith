---
name: db-migration
description: >
  Creates SQL migration files for the go-ecommerce-ddd-hexagonal-modular-monolith
  project. Covers migration naming, up/down files, column naming conventions,
  and versioning. Use when adding tables, modifying schemas, or running migrations.
metadata:
  author: ryvasa
  version: "1.0.0"
  domain: database
  triggers: migration, SQL, table, column, schema, ALTER, CREATE, PostgreSQL
  role: specialist
  scope: infrastructure
  output-format: sql
  related-skills: gorm-persistence
---

# Database Migration — Schema Management

## Migration Location

All migrations live in `migrations/` at project root:

```
migrations/
  000001_create_table_users.up.sql
  000001_create_table_users.down.sql
  000002_create_table_carts.up.sql
  000002_create_table_carts.down.sql
```

## File Naming Convention

```
{sequence}_{description}.{direction}.sql
```

- **sequence**: 6-digit zero-padded number (e.g., `000001`)
- **description**: snake_case describing the change
- **direction**: `up` for apply, `down` for rollback

## Up Migration Example

```sql
-- 000001_create_table_users.up.sql
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    email_verified_at TIMESTAMP,
    password VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    phone VARCHAR(255),
    birth_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);
```

## Down Migration Example

```sql
-- 000001_create_table_users.down.sql
DROP TABLE IF EXISTS users;
```

## Running Migrations

```bash
go run cmd/migrate/main.go
```

## Column Naming

- Use snake_case for column names
- Column names MUST reflect domain language
- Use `VARCHAR(36)` for UUID primary keys
- Use `TIMESTAMP` for datetime fields
- Add `created_at`, `updated_at`, `deleted_at` for GORM soft delete

## Constraints

### MUST DO
- Always create both up and down migration files
- Use `IF NOT EXISTS` / `IF EXISTS` for safety
- Version migrations sequentially
- Column names match persistence model GORM tags

### MUST NOT DO
- DO NOT use auto-migrate in production
- DO NOT modify existing applied migrations
- DO NOT skip down migration file
