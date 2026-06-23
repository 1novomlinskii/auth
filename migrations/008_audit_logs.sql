-- +goose Up

CREATE OR REPLACE FUNCTION set_user(p_user_id UUID)
RETURNS BOOLEAN AS $$
BEGIN
    PERFORM set_config('app.current_user_id', p_user_id::TEXT, true);
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE audit_logs (
    audit_id      UUID PRIMARY KEY DEFAULT uuidv7(),
    actor_id      UUID,
    action        VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id   UUID,
    tenant_id     UUID NOT NULL,
    ip_address    INET,
    result        VARCHAR(10) NOT NULL DEFAULT 'success',
    metadata      JSONB,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION fn_audit_trigger()
RETURNS TRIGGER AS $$
DECLARE
    v_actor_id   UUID;
    v_resource_id UUID;
    v_tenant_id  UUID;
BEGIN
    BEGIN
        v_actor_id := current_setting('app.current_user_id')::UUID;
    EXCEPTION WHEN OTHERS THEN
        v_actor_id := NULL;
    END;

    IF TG_OP = 'DELETE' THEN
        EXECUTE format('SELECT ($1).%I, ($1).%I', TG_ARGV[0], TG_ARGV[1])
        INTO v_resource_id, v_tenant_id USING OLD;
    ELSE
        EXECUTE format('SELECT ($1).%I, ($1).%I', TG_ARGV[0], TG_ARGV[1])
        INTO v_resource_id, v_tenant_id USING NEW;
    END IF;

    INSERT INTO audit_logs (actor_id, action, resource_type, resource_id, tenant_id)
    VALUES (
        v_actor_id,
        TG_TABLE_NAME || '.' || CASE TG_OP
            WHEN 'INSERT' THEN 'created'
            WHEN 'UPDATE' THEN 'updated'
            WHEN 'DELETE' THEN 'deleted'
        END,
        TG_TABLE_NAME,
        v_resource_id,
        COALESCE(v_tenant_id, '00000000-0000-0000-0000-000000000001')
    );
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE TRIGGER trg_users_audit AFTER INSERT OR UPDATE OR DELETE ON users
    FOR EACH ROW EXECUTE FUNCTION fn_audit_trigger('user_id', 'tenant_id');
CREATE TRIGGER trg_user_identities_audit AFTER INSERT OR UPDATE OR DELETE ON user_identities
    FOR EACH ROW EXECUTE FUNCTION fn_audit_trigger('identity_id', NULL);
CREATE TRIGGER trg_user_password_audit AFTER INSERT OR UPDATE OR DELETE ON user_password
    FOR EACH ROW EXECUTE FUNCTION fn_audit_trigger('user_id', NULL);
CREATE TRIGGER trg_user_sessions_audit AFTER INSERT OR UPDATE OR DELETE ON user_sessions
    FOR EACH ROW EXECUTE FUNCTION fn_audit_trigger('session_id', NULL);
CREATE TRIGGER trg_api_keys_audit AFTER INSERT OR UPDATE OR DELETE ON api_keys
    FOR EACH ROW EXECUTE FUNCTION fn_audit_trigger('api_key_id', 'tenant_id');
CREATE TRIGGER trg_roles_audit AFTER INSERT OR UPDATE OR DELETE ON roles
    FOR EACH ROW EXECUTE FUNCTION fn_audit_trigger('role_id', 'tenant_id');
CREATE TRIGGER trg_user_roles_audit AFTER INSERT OR UPDATE OR DELETE ON user_roles
    FOR EACH ROW EXECUTE FUNCTION fn_audit_trigger('user_id', NULL);

-- +goose Down
DROP TRIGGER IF EXISTS trg_users_audit ON users;
DROP TRIGGER IF EXISTS trg_user_identities_audit ON user_identities;
DROP TRIGGER IF EXISTS trg_user_password_audit ON user_password;
DROP TRIGGER IF EXISTS trg_user_sessions_audit ON user_sessions;
DROP TRIGGER IF EXISTS trg_api_keys_audit ON api_keys;
DROP TRIGGER IF EXISTS trg_roles_audit ON roles;
DROP TRIGGER IF EXISTS trg_user_roles_audit ON user_roles;
DROP FUNCTION IF EXISTS fn_audit_trigger();
DROP TABLE IF EXISTS audit_logs;
DROP FUNCTION IF EXISTS set_user(UUID);