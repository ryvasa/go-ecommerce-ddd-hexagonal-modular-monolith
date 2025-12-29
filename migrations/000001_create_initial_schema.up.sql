-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) COLLATE "pg_catalog"."default" PRIMARY KEY,
    email VARCHAR(255) COLLATE "pg_catalog"."default" NOT NULL UNIQUE,
    password VARCHAR(255) COLLATE "pg_catalog"."default" NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(36) COLLATE "pg_catalog"."default" PRIMARY KEY,
    user_id VARCHAR(36) COLLATE "pg_catalog"."default" NOT NULL,
    title VARCHAR(255) COLLATE "pg_catalog"."default" NOT NULL,
    done BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign key constraint separately
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_tasks_user_id'
    ) THEN
        ALTER TABLE tasks ADD CONSTRAINT fk_tasks_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END$$;

CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);

-- Create carts table
CREATE TABLE IF NOT EXISTS carts (
    id VARCHAR(36) COLLATE "pg_catalog"."default" PRIMARY KEY,
    user_id VARCHAR(36) COLLATE "pg_catalog"."default" NOT NULL,
    title VARCHAR(255) COLLATE "pg_catalog"."default" NOT NULL,
    done BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign key constraint separately
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_carts_user_id'
    ) THEN
        ALTER TABLE carts ADD CONSTRAINT fk_carts_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END$$;

CREATE INDEX IF NOT EXISTS idx_carts_user_id ON carts(user_id);

-- Create casbin_rule table for authorization
CREATE TABLE IF NOT EXISTS casbin_rule (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(100) COLLATE "pg_catalog"."default",
    v0 VARCHAR(100) COLLATE "pg_catalog"."default",
    v1 VARCHAR(100) COLLATE "pg_catalog"."default",
    v2 VARCHAR(100) COLLATE "pg_catalog"."default",
    v3 VARCHAR(100) COLLATE "pg_catalog"."default",
    v4 VARCHAR(100) COLLATE "pg_catalog"."default",
    v5 VARCHAR(100) COLLATE "pg_catalog"."default"
);

-- Add unique constraint if not exists
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'unique_key_sqlx_migrator'
    ) THEN
        ALTER TABLE casbin_rule ADD CONSTRAINT unique_key_sqlx_migrator 
        UNIQUE(ptype, v0, v1, v2, v3, v4, v5);
    END IF;
END$$;

CREATE INDEX IF NOT EXISTS idx_casbin_rule ON casbin_rule(ptype, v0, v1);
