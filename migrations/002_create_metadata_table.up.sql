BEGIN;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS metadata (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    is_file BOOLEAN NOT NULL,
    public BOOLEAN NOT NULL,
    mime TEXT NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    owner_id UUID NOT NULL,
    json_data JSON NOT NULL,
    file_size BIGINT,
    url TEXT,
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_metadata_owner ON metadata(owner_id);

COMMIT;
