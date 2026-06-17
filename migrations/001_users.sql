-- +goose Up
CREATE TABLE users (
    user_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000001',
    global_handle VARCHAR(50) UNIQUE NOT NULL,
    display_name  VARCHAR(255),
    email_primary VARCHAR(255) UNIQUE NOT NULL,
    locale        VARCHAR(10),
    timezone      VARCHAR(50),
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ,
    is_active     BOOLEAN DEFAULT TRUE
);

-- +goose Down
DROP TABLE IF EXISTS users;