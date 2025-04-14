-- +goose Up
-- +goose StatementBegin
ALTER TABLE links
    ADD COLUMN owner_id UUID REFERENCES users(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE links
    DROP CONSTRAINT IF EXISTS fk_links_users,
    DROP COLUMN IF EXISTS owner_id;
-- +goose StatementEnd