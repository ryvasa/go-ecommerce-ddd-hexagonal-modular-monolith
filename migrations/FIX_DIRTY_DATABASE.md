# Fix Dirty Database on Railway

Jika kamu mengalami error **"Dirty database version 1"** di Railway, ikuti langkah ini:

## Option 1: Reset Database via Railway Dashboard (Recommended)

1. **Buka Railway Dashboard** → Your Project → PostgreSQL Service
2. **Klik tab "Data"**
3. **Run SQL Query:**

```sql
-- Check dirty state
SELECT * FROM schema_migrations;

-- Drop schema_migrations table to reset
DROP TABLE IF EXISTS schema_migrations;

-- Optionally, drop all tables to start fresh
DROP TABLE IF EXISTS casbin_rule;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS users;
```

4. **Redeploy** aplikasi dari Railway dashboard
5. Migration akan jalan ulang dari fresh state

## Option 2: Force Migration Version (Quick Fix)

Jika kamu perlu quick fix tanpa reset database:

**Via Railway's PostgreSQL Console:**

```sql
-- Set migration version to 0 (clean state)
DELETE FROM schema_migrations;
```

Lalu redeploy aplikasi.

## Option 3: Manual Migration via Railway CLI

```bash
# Install Railway CLI
npm i -g @railway/cli

# Login
railway login

# Link to project
railway link

# Connect to database and run SQL
railway run psql $DATABASE_URL

# In psql console:
DROP TABLE IF EXISTS schema_migrations;
\q

# Redeploy
git commit --allow-empty -m "Trigger redeploy"
git push
```

## Prevention

Untuk mencegah dirty state di future updates:

1. **Test migrations locally first:**

   ```bash
   # Local test
   go run cmd/server/main.go
   ```

2. **Gunakan Railway Preview Environments** untuk test sebelum merge ke main

3. **Backup database** sebelum deploy breaking changes

## What Happened?

Error "Dirty database version 1" terjadi karena:

- Migration started
- Migration failed di tengah (foreign key error)
- Database marked as "dirty" (stuck state)
- Subsequent attempts can't proceed

Solusinya: reset migration state atau force version.
