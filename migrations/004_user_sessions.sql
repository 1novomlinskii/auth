-- +goose Up
CREATE TABLE user_sessions (
    session_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    identity_id   UUID REFERENCES user_identities(identity_id) ON DELETE SET NULL,
    refresh_token TEXT UNIQUE NOT NULL,
    device_info   TEXT,
    ip_address    INET,
    expires_at    TIMESTAMPTZ NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    revoked_at    TIMESTAMPTZ
);
CREATE INDEX idx_sessions_user ON user_sessions (user_id) WHERE revoked_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS user_sessions;