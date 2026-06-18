-- +goose Up
CREATE TABLE user_password_hash (
    user_id             UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
    password_hash       TEXT NOT NULL,
    reset_token_hash    TEXT,
    reset_token_expires_at TIMESTAMPTZ,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS user_password_hash;