CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- URLs table
CREATE TABLE IF NOT EXISTS urls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    original_url TEXT NOT NULL,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    click_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Index so redirect lookup is fast
CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);