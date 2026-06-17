-- +goose Up
CREATE TABLE user_identities (
    identity_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    provider         VARCHAR(20) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    provider_email   VARCHAR(255),
    provider_name    VARCHAR(255),
    avatar_url       TEXT,
    last_login_at    TIMESTAMPTZ,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (provider, provider_user_id)
);

-- +goose Down
DROP TABLE IF EXISTS user_identities;