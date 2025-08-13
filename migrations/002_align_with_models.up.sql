-- Migration to align database schema with Go models
-- This migration updates the existing schema to match the defined models

-- Update users table to match User model
ALTER TABLE users
DROP COLUMN IF EXISTS username,
DROP COLUMN IF EXISTS bio,
DROP COLUMN IF EXISTS avatar,
DROP COLUMN IF EXISTS email_verified,
DROP COLUMN IF EXISTS last_login,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE users
ADD COLUMN IF NOT EXISTS name VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN IF NOT EXISTS role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (
    role IN ('admin', 'user', 'guest')
);

-- Update questions table to match Question model
ALTER TABLE questions
DROP COLUMN IF EXISTS title,
DROP COLUMN IF EXISTS category_id,
DROP COLUMN IF EXISTS is_anonymous,
DROP COLUMN IF EXISTS view_count,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE questions
ADD COLUMN IF NOT EXISTS like_count INTEGER DEFAULT 0;

-- Update answers table to match Answer model
ALTER TABLE answers
DROP COLUMN IF EXISTS is_accepted,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE answers RENAME COLUMN content TO answer;

-- Add question_id column if it doesn't exist (it should exist from previous migration)
-- This is already correct in the current schema

-- Create rooms table to match Room model
CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Drop tables that are not represented in models
DROP TABLE IF EXISTS categories CASCADE;

DROP TABLE IF EXISTS tags CASCADE;

DROP TABLE IF EXISTS question_tags CASCADE;

DROP TABLE IF EXISTS votes CASCADE;

-- Update indexes to match new structure
DROP INDEX IF EXISTS idx_questions_category_id;

DROP INDEX IF EXISTS idx_questions_created_at;

DROP INDEX IF EXISTS idx_votes_target;

DROP INDEX IF EXISTS idx_votes_user_id;

-- Add new indexes for performance
CREATE INDEX IF NOT EXISTS idx_questions_like_count ON questions (like_count);

CREATE INDEX IF NOT EXISTS idx_rooms_owner_id ON rooms (owner_id);

-- Keep existing indexes that are still relevant
-- idx_questions_user_id, idx_answers_question_id, idx_answers_user_id are still valid