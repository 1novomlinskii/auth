-- +goose Up
CREATE TABLE roles (
    role_id       UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id     UUID NOT NULL,
    name          TEXT NOT NULL,
    permissions   TEXT[] NOT NULL DEFAULT '{}',
    is_system     BOOLEAN DEFAULT FALSE,
    description   TEXT,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ,
    UNIQUE (tenant_id, name)
);

-- +goose Down
DROP TABLE IF EXISTS roles;