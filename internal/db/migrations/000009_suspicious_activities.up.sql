CREATE TABLE suspicious_activities (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users (id) ON DELETE SET NULL,
    activity_type VARCHAR(100) NOT NULL,
    ip_address INET NOT NULL,
    user_agent TEXT,
    description TEXT,
    metadata JSONB,
    severity VARCHAR(20) DEFAULT 'medium',
    resolved BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    CHECK (severity IN ('low', 'medium', 'high', 'critical'))
);


