-- Create user_reactions table
-- +migrate up
CREATE TABLE user_reactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES user_sessions(id) ON DELETE CASCADE,
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    reaction_type VARCHAR(50) NOT NULL, -- 'like', 'dislike', etc.
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

-- Prevent duplicate reactions from same user to same message
UNIQUE(session_id, message_id, reaction_type) );

-- Create room_creators table to track who created each room
CREATE TABLE room_creators (
    room_id UUID PRIMARY KEY REFERENCES rooms (id) ON DELETE CASCADE,
    creator_session_id UUID NOT NULL REFERENCES user_sessions (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_user_reactions_session ON user_reactions (session_id);

CREATE INDEX idx_user_reactions_room ON user_reactions (room_id);

CREATE INDEX idx_user_reactions_message ON user_reactions (message_id);

CREATE INDEX idx_user_reactions_type ON user_reactions (reaction_type);

CREATE INDEX idx_room_creators_session ON room_creators (creator_session_id);

-- +migrate down
DROP TABLE room_creators;

DROP TABLE user_reactions;