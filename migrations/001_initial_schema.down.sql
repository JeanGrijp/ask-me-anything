DROP INDEX IF EXISTS idx_magic_links_expires_at;

DROP INDEX IF EXISTS idx_magic_links_email;

DROP INDEX IF EXISTS idx_magic_links_token;

DROP INDEX IF EXISTS idx_votes_user_id;

DROP INDEX IF EXISTS idx_votes_target;

DROP INDEX IF EXISTS idx_answers_user_id;

DROP INDEX IF EXISTS idx_answers_question_id;

DROP INDEX IF EXISTS idx_questions_created_at;

DROP INDEX IF EXISTS idx_questions_category_id;

DROP INDEX IF EXISTS idx_questions_user_id;

DROP TABLE IF EXISTS magic_links;

DROP TABLE IF EXISTS votes;

DROP TABLE IF EXISTS question_tags;

DROP TABLE IF EXISTS tags;

DROP TABLE IF EXISTS answers;

DROP TABLE IF EXISTS questions;

DROP TABLE IF EXISTS categories;

DROP TABLE IF EXISTS users;