-- Add is_answered flag to questions table
-- This flag allows room owners to mark questions as answered/unanswered

ALTER TABLE questions ADD COLUMN is_answered BOOLEAN DEFAULT FALSE;

-- Add index for better performance when filtering answered/unanswered questions
CREATE INDEX idx_questions_is_answered ON questions (is_answered);

-- Add comment to document the purpose
COMMENT ON COLUMN questions.is_answered IS 'Flag to mark if the question has been answered by the room owner';