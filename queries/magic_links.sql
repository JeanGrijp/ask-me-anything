-- name: CreateMagicLink :one
INSERT INTO
    magic_links (token, email, expires_at)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetMagicLink :one
SELECT *
FROM magic_links
WHERE
    token = $1
    AND used = FALSE
    AND expires_at > NOW()
LIMIT 1;

-- name: UseMagicLink :exec
UPDATE magic_links
SET
    used = TRUE,
    used_at = CURRENT_TIMESTAMP
WHERE
    token = $1;

-- name: DeleteExpiredMagicLinks :exec
DELETE FROM magic_links WHERE expires_at < NOW() OR used = TRUE;

-- name: GetMagicLinkByEmail :many
SELECT *
FROM magic_links
WHERE
    email = $1
    AND used = FALSE
    AND expires_at > NOW()
ORDER BY created_at DESC;