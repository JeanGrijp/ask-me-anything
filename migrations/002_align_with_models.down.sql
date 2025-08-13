-- Rollback migration 002 - restore original schema structure

-- Restore rooms table removal
DROP TABLE IF EXISTS rooms;

-- Restore votes table
CREATE TABLE IF NOT EXISTS votes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    target_id INTEGER NOT NULL,
    target_type VARCHAR(20) NOT NULL CHECK (
        target_type IN ('question', 'answer')
    ),
    vote_type VARCHAR(10) NOT NULL CHECK (vote_type IN ('up', 'down')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (
        user_id,
        target_id,
        target_type
    )
);

-- Restore tags tables
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS question_tags (
    question_id INTEGER NOT NULL REFERENCES questions (id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (question_id, tag_id)
);

-- Restore categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    color VARCHAR(7)
);

-- Restore answers table structure
ALTER TABLE answers RENAME COLUMN answer TO content;

ALTER TABLE answers
ADD COLUMN IF NOT EXISTS is_accepted BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Restore questions table structure
ALTER TABLE questions DROP COLUMN IF EXISTS like_count;

ALTER TABLE questions
ADD COLUMN IF NOT EXISTS title VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories (id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS is_anonymous BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS view_count INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Restore users table structure
ALTER TABLE users
DROP COLUMN IF EXISTS name,
DROP COLUMN IF EXISTS role;

ALTER TABLE users
ADD COLUMN IF NOT EXISTS username VARCHAR(255) UNIQUE,
ADD COLUMN IF NOT EXISTS bio TEXT,
ADD COLUMN IF NOT EXISTS avatar VARCHAR(255),
ADD COLUMN IF NOT EXISTS email_verified BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS last_login TIMESTAMP,
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Restore indexes
DROP INDEX IF EXISTS idx_questions_like_count;

DROP INDEX IF EXISTS idx_rooms_owner_id;

CREATE INDEX IF NOT EXISTS idx_questions_category_id ON questions (category_id);

CREATE INDEX IF NOT EXISTS idx_questions_created_at ON questions (created_at);

CREATE INDEX IF NOT EXISTS idx_votes_target ON votes (target_id, target_type);

CREATE INDEX IF NOT EXISTS idx_votes_user_id ON votes (user_id);