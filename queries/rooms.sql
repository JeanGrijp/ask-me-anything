-- Room queries to match Room model
-- name: GetRoom :one
SELECT * FROM rooms WHERE id = $1 LIMIT 1;

-- name: ListRooms :many
SELECT * FROM rooms ORDER BY created_at DESC;

-- name: ListRoomsByOwner :many
SELECT * FROM rooms WHERE owner_id = $1 ORDER BY created_at DESC;

-- name: CreateRoom :one
INSERT INTO rooms (name, owner_id) VALUES ($1, $2) RETURNING *;

-- name: UpdateRoom :one
UPDATE rooms SET name = $2 WHERE id = $1 RETURNING *;

-- name: DeleteRoom :exec
DELETE FROM rooms WHERE id = $1;