-- +migrate Up
-- wordbook_entries table
CREATE TABLE IF NOT EXISTS wordbook_entries (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL DEFAULT 'default',
    word VARCHAR(128) NOT NULL,
    short_definition TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, word)
);
CREATE INDEX IF NOT EXISTS idx_wordbook_user_created ON wordbook_entries(user_id, created_at DESC);

-- dictionary_cache table
CREATE TABLE IF NOT EXISTS dictionary_cache (
    word VARCHAR(128) PRIMARY KEY,
    data JSONB NOT NULL,
    source VARCHAR(64) NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_cache_expires ON dictionary_cache(expires_at);

-- +migrate Down
DROP TABLE IF EXISTS dictionary_cache;
DROP TABLE IF EXISTS wordbook_entries;
