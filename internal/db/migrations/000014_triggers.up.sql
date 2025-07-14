CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_profiles_updated_at BEFORE UPDATE ON user_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_oauth_accounts_updated_at BEFORE UPDATE ON oauth_accounts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


CREATE OR REPLACE FUNCTION log_password_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.password_hash != NEW.password_hash THEN
        INSERT INTO password_history (user_id, password_hash)
        VALUES (NEW.id, OLD.password_hash);

        INSERT INTO audit_logs (user_id, action, resource_type, resource_id, details)
        VALUES (NEW.id, 'password_change', 'user', NEW.id, '{"method": "direct_change"}'::jsonb);
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER log_user_password_change AFTER UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION log_password_change();

CREATE OR REPLACE FUNCTION cleanup_expired_data()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER := 0;
BEGIN
    DELETE FROM refresh_tokens WHERE expires_at < NOW();
    GET DIAGNOSTICS deleted_count = ROW_COUNT;

    DELETE FROM email_verifications WHERE expires_at < NOW();

    DELETE FROM password_resets WHERE expires_at < NOW();

    DELETE FROM account_lockouts WHERE expires_at < NOW();

    DELETE FROM login_attempts WHERE created_at < NOW() - INTERVAL '90 days';

    DELETE FROM password_history
    WHERE id NOT IN (
        SELECT id FROM (
            SELECT id, ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at DESC) as rn
            FROM password_history
        ) ranked WHERE rn <= 10
    );

    RETURN deleted_count;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION check_active_lockout()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM account_lockouts
        WHERE user_id = NEW.user_id
          AND lockout_type = NEW.lockout_type
          AND expires_at > NOW()
          AND (TG_OP = 'INSERT' OR id <> NEW.id)
    ) THEN
        RAISE EXCEPTION 'An active lockout of type "%" already exists for this user.', NEW.lockout_type;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_active_lockout BEFORE INSERT OR UPDATE ON account_lockouts FOR EACH ROW EXECUTE FUNCTION check_active_lockout();
