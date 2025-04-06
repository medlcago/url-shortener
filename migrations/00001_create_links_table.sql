-- +goose Up
-- +goose StatementBegin
CREATE TABLE links
(
    id           UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    original_url TEXT                NOT NULL,
    alias        VARCHAR(255) UNIQUE NOT NULL,
    expires_at   TIMESTAMP WITH TIME ZONE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC')
);
CREATE INDEX IF NOT EXISTS idx_alias ON links (alias);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_alias;
DROP TABLE IF EXISTS links;
-- +goose StatementEnd