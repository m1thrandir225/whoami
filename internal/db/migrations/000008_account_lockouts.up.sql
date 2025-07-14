CREATE TABLE account_lockouts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    ip_address INET,
    lockout_type VARCHAR(50) NOT NULL CHECK (lockout_type IN ('account', 'ip', 'both')),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
