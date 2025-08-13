-- Answers queries aligned with Answer model

-- name: GetAnswer :one
SELECT * FROM answers WHERE id = $1 LIMIT 1;

-- name: ListAnswersByQuestion :many
SELECT * FROM answers WHERE question_id = $1 ORDER BY created_at ASC;

-- name: ListAnswersByUser :many
SELECT * FROM answers WHERE user_id = $1 ORDER BY created_at DESC;

-- name: CreateAnswer :one
INSERT INTO
    answers (answer, question_id, user_id)
VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateAnswer :one
UPDATE answers SET answer = $2 WHERE id = $1 RETURNING *;

-- name: DeleteAnswer :exec
DELETE FROM answers WHERE id = $1;