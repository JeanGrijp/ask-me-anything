-- Create user_sessions table
-- +migrate up
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_token VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    last_activity TIMESTAMP NOT NULL DEFAULT NOW(),

-- Optional user info that can be collected over time
username VARCHAR(100), email VARCHAR(255),

-- Session metadata
user_agent TEXT, ip_address INET );

-- Create indexes for performance
CREATE INDEX idx_user_sessions_token ON user_sessions (session_token);

CREATE INDEX idx_user_sessions_expires ON user_sessions (expires_at);

CREATE INDEX idx_user_sessions_last_activity ON user_sessions (last_activity);

-- +migrate down
DROP TABLE user_sessions;