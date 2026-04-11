---
description: create a new database migration (up + down SQL files)
---

# New Migration — Database Schema Change

Workflow to create a new SQL migration for schema changes.

---

## 1. Determine Next Sequence

// turbo

```bash
ls -la migrations/*.up.sql | tail -1
```

Get the last migration number and increment by 1.

---

## 2. Create Migration Files

Create both up and down migration files:

```bash
NEXT=000002
DESC=create_table_orders
touch migrations/${NEXT}_${DESC}.up.sql
touch migrations/${NEXT}_${DESC}.down.sql
```

---

## 3. Write Up Migration

Write the SQL for applying the change. Follow conventions:

- Use `snake_case` for table and column names
- Use `VARCHAR(36)` for UUID primary keys
- Add `created_at`, `updated_at`, `deleted_at` columns
- Use `IF NOT EXISTS` for safety

---

## 4. Write Down Migration

Write the SQL for rolling back the change:

- Use `DROP TABLE IF EXISTS` for table removals
- For ALTER TABLE changes, write the reverse operation

---

## 5. Test Migration

```bash
go run cmd/migrate/main.go
```

Verify migration applies successfully.

---

## 6. Verify Model Alignment

Ensure the GORM persistence model matches the migration schema:
- Column names match
- Types match
- Constraints match
