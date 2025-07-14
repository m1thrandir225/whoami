DROP TRIGGER IF EXISTS trigger_check_active_lockout ON account_lockouts;
DROP FUNCTION IF EXISTS check_active_lockout();
DROP FUNCTION IF EXISTS cleanup_expired_data();
DROP TRIGGER IF EXISTS log_user_password_change ON users;
DROP FUNCTION IF EXISTS log_password_change();
DROP TRIGGER IF EXISTS update_oauth_accounts_updated_at ON oauth_accounts;
DROP TRIGGER IF EXISTS update_user_profiles_updated_at ON user_profiles;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
