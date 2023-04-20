CREATE TABLE IF NOT EXISTS "media"
(
    hash        BYTEA PRIMARY KEY,
    path        TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL,
    imported_at TIMESTAMPTZ NOT NULL DEFAULT now()
);