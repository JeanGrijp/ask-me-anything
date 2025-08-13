CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    bio TEXT,
    avatar VARCHAR(255),
    email_verified BOOLEAN DEFAULT FALSE,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS magic_links (
    id SERIAL PRIMARY KEY,
    token VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    color VARCHAR(7)
);

CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories (id) ON DELETE SET NULL,
    is_anonymous BOOLEAN DEFAULT FALSE,
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS answers (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    question_id INTEGER NOT NULL REFERENCES questions (id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    is_accepted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS question_tags (
    question_id INTEGER NOT NULL REFERENCES questions (id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (question_id, tag_id)
);

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

-- Indexes for better performance
CREATE INDEX idx_questions_user_id ON questions (user_id);

CREATE INDEX idx_questions_category_id ON questions (category_id);

CREATE INDEX idx_questions_created_at ON questions (created_at);

CREATE INDEX idx_answers_question_id ON answers (question_id);

CREATE INDEX idx_answers_user_id ON answers (user_id);

CREATE INDEX idx_votes_target ON votes (target_id, target_type);

CREATE INDEX idx_votes_user_id ON votes (user_id);

CREATE INDEX idx_magic_links_token ON magic_links (token);

CREATE INDEX idx_magic_links_email ON magic_links (email);

CREATE INDEX idx_magic_links_expires_at ON magic_links (expires_at);