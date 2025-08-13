-- Remove is_answered flag from questions table

-- Remove the index first
DROP INDEX IF EXISTS idx_questions_is_answered;

-- Remove the column
ALTER TABLE questions DROP COLUMN IF EXISTS is_answered;