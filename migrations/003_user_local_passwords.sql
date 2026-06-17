-- +goose Up
CREATE TABLE user_local_passwords (
    user_id       UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS user_local_passwords;