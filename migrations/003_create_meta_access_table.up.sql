BEGIN;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS meta_access (
    meta_id UUID NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY (meta_id, user_id),
    FOREIGN KEY (meta_id) REFERENCES metadata(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

COMMIT;
