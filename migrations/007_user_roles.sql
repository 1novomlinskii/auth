-- +goose Up
CREATE TABLE user_roles (
    user_id       UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    role_id       UUID NOT NULL REFERENCES roles(role_id) ON DELETE CASCADE,
    assigned_at   TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

-- +goose Down
DROP TABLE IF EXISTS user_roles;