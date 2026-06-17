-- +goose Up
CREATE TABLE api_keys (
    api_key_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL,
    service_name  TEXT NOT NULL,
    key_prefix    VARCHAR(8) NOT NULL,
    key_hash      TEXT NOT NULL,
    permissions   TEXT[] NOT NULL DEFAULT '{}',
    created_by    UUID REFERENCES users(user_id) ON DELETE SET NULL,
    expires_at    TIMESTAMPTZ,
    is_active     BOOLEAN DEFAULT TRUE,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (key_hash)
);

-- +goose Down
DROP TABLE IF EXISTS api_keys;