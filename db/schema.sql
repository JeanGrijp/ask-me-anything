-- Final schema for Ask Me Anything Application
-- This reflects the database structure after all migrations are applied

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL DEFAULT '',
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (
        role IN ('admin', 'user', 'guest')
    ),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Rooms table
CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Questions table
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    like_count INTEGER DEFAULT 0,
    is_answered BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Answers table
CREATE TABLE answers (
    id SERIAL PRIMARY KEY,
    answer TEXT NOT NULL,
    question_id INTEGER NOT NULL REFERENCES questions (id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Magic links table for authentication
CREATE TABLE magic_links (
    id SERIAL PRIMARY KEY,
    token VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT false,
    used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_questions_user_id ON questions (user_id);

CREATE INDEX idx_questions_like_count ON questions (like_count);

CREATE INDEX idx_questions_is_answered ON questions (is_answered);

CREATE INDEX idx_answers_question_id ON answers (question_id);

CREATE INDEX idx_answers_user_id ON answers (user_id);

CREATE INDEX idx_rooms_owner_id ON rooms (owner_id);

CREATE INDEX idx_magic_links_token ON magic_links (token);

CREATE INDEX idx_magic_links_email ON magic_links (email);

CREATE INDEX idx_magic_links_expires_at ON magic_links (expires_at);