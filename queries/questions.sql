-- name: GetQuestion :one
SELECT * FROM questions WHERE id = $1 LIMIT 1;

-- name: ListQuestions :many
SELECT * FROM questions ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListQuestionsByUser :many
SELECT * FROM questions WHERE user_id = $1 ORDER BY created_at DESC;

-- Updated questions queries to match Question model
-- name: GetQuestion :one
SELECT * FROM questions WHERE id = $1 LIMIT 1;

-- name: ListQuestions :many
SELECT *
FROM questions
ORDER BY like_count DESC, created_at DESC
LIMIT $1
OFFSET
    $2;

-- name: ListQuestionsByUser :many
SELECT * FROM questions WHERE user_id = $1 ORDER BY created_at DESC;

-- name: CreateQuestion :one
INSERT INTO
    questions (content, user_id, like_count)
VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateQuestion :one
UPDATE questions
SET
    content = $2,
    like_count = $3
WHERE
    id = $1 RETURNING *;

-- name: IncrementQuestionLikes :exec
UPDATE questions SET like_count = like_count + 1 WHERE id = $1;

-- name: DecrementQuestionLikes :exec
UPDATE questions
SET
    like_count = like_count - 1
WHERE
    id = $1
    AND like_count > 0;

-- name: DeleteQuestion :exec
DELETE FROM questions WHERE id = $1;

-- name: CreateQuestion :one
INSERT INTO
    questions (
        title,
        content,
        user_id,
        category_id,
        is_anonymous
    )
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateQuestion :one
UPDATE questions
SET
    title = $2,
    content = $3,
    category_id = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1 RETURNING *;

-- name: IncrementQuestionViews :exec
UPDATE questions SET view_count = view_count + 1 WHERE id = $1;

-- name: DeleteQuestion :exec
DELETE FROM questions WHERE id = $1;