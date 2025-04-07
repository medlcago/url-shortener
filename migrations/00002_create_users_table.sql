-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id                UUID PRIMARY KEY                  DEFAULT gen_random_uuid(),
    email             VARCHAR(255)             NOT NULL,
    password          VARCHAR(255)             NOT NULL,
    is_active         BOOLEAN                  NOT NULL DEFAULT TRUE,
    is_email_verified BOOLEAN                  NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    login_date        TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT email_unique UNIQUE (email)
);

CREATE OR REPLACE FUNCTION update_updated_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = (NOW() AT TIME ZONE 'UTC');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_updated_at
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_update_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at();
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;
-- +goose StatementEnd
