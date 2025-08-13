-- Questions queries aligned with Question model

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
VALUES ($1, $2, COALESCE($3, 0)) RETURNING *;

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