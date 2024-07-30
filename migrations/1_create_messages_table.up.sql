CREATE TABLE IF NOT EXISTS messages
(
    id           SERIAL PRIMARY KEY,
    content      TEXT        NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed    BOOLEAN     NOT NULL DEFAULT FALSE,
    processed_at TIMESTAMPTZ          DEFAULT NULL
);