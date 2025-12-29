# Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema management.

## Auto-Migration on Startup

Migrations run automatically when the application starts. This works for both local development and Railway deployment.

## Migration Files Location

All migration files are in the `migrations/` directory:

```
migrations/
├── 000001_create_initial_schema.up.sql    # Creates tables
└── 000001_create_initial_schema.down.sql  # Drops tables
```

## Creating New Migrations

### Option 1: Manual (Recommended)

Create two files with sequential numbering:

```bash
# Up migration
touch migrations/000002_add_user_role.up.sql

# Down migration
touch migrations/000002_add_user_role.down.sql
```

### Option 2: Using migrate CLI (Optional)

```bash
# Install CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migration files
migrate create -ext sql -dir migrations -seq migration_name
```

## Example Migration

**migrations/000002_add_user_role.up.sql:**

```sql
ALTER TABLE users ADD COLUMN role VARCHAR(50) DEFAULT 'user';
CREATE INDEX idx_users_role ON users(role);
```

**migrations/000002_add_user_role.down.sql:**

```sql
DROP INDEX IF EXISTS idx_users_role;
ALTER TABLE users DROP COLUMN IF EXISTS role;
```

## Running Migrations

### Automatic (on app startup)

Migrations run automatically when you start the server:

```bash
go run cmd/server/main.go
```

Output:

```
🔄 Running database migrations...
✅ Migrations completed successfully
```

### Manual (for testing)

You can also run migrations manually using the CLI:

```bash
# Up
migrate -path migrations -database $POSTGRES_DSN up

# Down (rollback all)
migrate -path migrations -database $POSTGRES_DSN down

# Down (rollback 1 step)
migrate -path migrations -database $POSTGRES_DSN down 1

# Goto specific version
migrate -path migrations -database $POSTGRES_DSN goto 2

# Force version (if stuck)
migrate -path migrations -database $POSTGRES_DSN force 1
```

## Railway Deployment

On Railway, migrations run automatically because:

1. Railway provides `DATABASE_URL` environment variable
2. Our code auto-detects and uses this variable
3. Migrations execute before the server starts serving requests

**Deployment flow:**

```
git push → Railway builds → App starts → Migrations run → Server ready
```

## Local Development

For local development, set `POSTGRES_DSN` in your `.env` file:

```env
POSTGRES_DSN=postgres://user:password@localhost:5432/dbname?sslmode=disable
```

## Troubleshooting

### Migration fails with "dirty database"

```bash
# Check current version
migrate -path migrations -database $POSTGRES_DSN version

# Force to a known good version
migrate -path migrations -database $POSTGRES_DSN force VERSION_NUMBER
```

### Schema already exists

The migrations use `IF NOT EXISTS` clauses, so they're safe to run multiple times.

### Need to rollback

```bash
# Rollback last migration
migrate -path migrations -database $POSTGRES_DSN down 1
```

## Best Practices

1. ✅ Always create both `.up.sql` and `.down.sql` files
2. ✅ Test migrations locally before deploying
3. ✅ Use `IF EXISTS` / `IF NOT EXISTS` for idempotency
4. ✅ Small, focused migrations (one change per migration)
5. ✅ Never modify existing migration files (create new ones instead)
6. ✅ Add indexes for foreign keys
7. ✅ Use transactions when possible

## Current Schema

- **users**: User authentication and profile
- **tasks**: User tasks
- **carts**: Shopping cart items
- **casbin_rule**: Authorization policies (Casbin)
