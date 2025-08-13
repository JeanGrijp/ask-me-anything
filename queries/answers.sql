-- name: GetAnswer :one
SELECT * FROM answers WHERE id = $1 LIMIT 1;

-- Updated answers queries to match Answer model
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

-- name: ListAnswersByUser :many
SELECT * FROM answers WHERE user_id = $1 ORDER BY created_at DESC;

-- name: CreateAnswer :one
INSERT INTO
    answers (content, question_id, user_id)
VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateAnswer :one
UPDATE answers
SET
    content = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1 RETURNING *;

-- name: AcceptAnswer :exec
UPDATE answers SET is_accepted = TRUE WHERE id = $1;

-- name: UnacceptAnswer :exec
UPDATE answers SET is_accepted = FALSE WHERE id = $1;

-- name: DeleteAnswer :exec
DELETE FROM answers WHERE id = $1;