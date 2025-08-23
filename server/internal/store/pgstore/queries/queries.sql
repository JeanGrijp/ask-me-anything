-- name: GetRoom :one
SELECT "id", "theme" FROM rooms WHERE id = $1;

-- name: GetRooms :many
SELECT "id", "theme" FROM rooms;

-- name: InsertRoom :one
INSERT INTO rooms ("theme") VALUES ($1) RETURNING "id";

-- name: GetMessage :one
SELECT "id", "room_id", "message", "reaction_count", "answered"
FROM messages
WHERE
    id = $1;

-- name: GetRoomMessages :many
SELECT "id", "room_id", "message", "reaction_count", "answered"
FROM messages
WHERE
    room_id = $1;

-- name: GetRoomMessagesWithUserReactions :many
SELECT
    m.id,
    m.room_id,
    m.message,
    m.reaction_count,
    m.answered,
    CASE
        WHEN ur.id IS NOT NULL THEN true
        ELSE false
    END as user_reacted
FROM messages m
    LEFT JOIN user_reactions ur ON (
        m.id = ur.message_id
        AND ur.session_id = (
            SELECT id
            FROM user_sessions
            WHERE
                session_token = $2
                AND expires_at > NOW()
        )
        AND ur.reaction_type = 'like'
    )
WHERE
    m.room_id = $1
ORDER BY m.id;

-- name: InsertMessage :one
INSERT INTO
    messages ("room_id", "message")
VALUES ($1, $2) RETURNING "id";

-- name: ReactToMessage :one
UPDATE messages
SET
    reaction_count = reaction_count + 1
WHERE
    id = $1 RETURNING reaction_count;

-- name: RemoveReactionFromMessage :one
UPDATE messages
SET
    reaction_count = reaction_count - 1
WHERE
    id = $1 RETURNING reaction_count;

-- name: MarkMessageAsAnswered :exec
UPDATE messages SET answered = true WHERE id = $1;

-- User Session Operations
-- name: CreateUserSession :one
INSERT INTO
    user_sessions (
        "session_token",
        "expires_at",
        "user_agent",
        "ip_address"
    )
VALUES ($1, $2, $3, $4) RETURNING "id",
    "session_token",
    "created_at",
    "expires_at";

-- name: GetUserSession :one
SELECT "id", "session_token", "created_at", "expires_at", "last_activity", "username", "email"
FROM user_sessions
WHERE
    session_token = $1
    AND expires_at > NOW();

-- name: UpdateSessionActivity :exec
UPDATE user_sessions
SET
    last_activity = NOW(),
    expires_at = $2
WHERE
    session_token = $1;

-- name: DeleteUserSession :exec
DELETE FROM user_sessions WHERE session_token = $1;

-- name: CleanExpiredSessions :exec
DELETE FROM user_sessions WHERE expires_at < NOW();

-- Room Creator Operations
-- name: SetRoomCreator :exec
INSERT INTO
    room_creators (
        "room_id",
        "creator_session_id"
    )
VALUES ($1, $2) ON CONFLICT (room_id) DO NOTHING;

-- name: GetRoomCreator :one
SELECT us.id, us.session_token, us.username
FROM
    room_creators rc
    JOIN user_sessions us ON rc.creator_session_id = us.id
WHERE
    rc.room_id = $1;

-- name: IsRoomCreator :one
SELECT EXISTS (
        SELECT 1
        FROM
            room_creators rc
            JOIN user_sessions us ON rc.creator_session_id = us.id
        WHERE
            rc.room_id = $1
            AND us.session_token = $2
    ) as is_creator;

-- name: GetUserRooms :many
SELECT r.id, r.theme, rc.created_at
FROM
    rooms r
    JOIN room_creators rc ON r.id = rc.room_id
    JOIN user_sessions us ON rc.creator_session_id = us.id
WHERE
    us.session_token = $1
ORDER BY rc.created_at DESC;

-- User Reaction Operations
-- name: AddUserReaction :exec
INSERT INTO
    user_reactions (
        "session_id",
        "room_id",
        "message_id",
        "reaction_type"
    )
VALUES ($1, $2, $3, $4) ON CONFLICT (
        session_id,
        message_id,
        reaction_type
    ) DO NOTHING;

-- name: RemoveUserReaction :exec
DELETE FROM user_reactions
WHERE
    session_id = $1
    AND message_id = $2
    AND reaction_type = $3;

-- name: GetUserReaction :one
SELECT "id", "reaction_type", "created_at"
FROM user_reactions
WHERE
    session_id = $1
    AND message_id = $2
    AND reaction_type = $3;

-- name: GetMessageReactions :many
SELECT ur.reaction_type, COUNT(*) as count
FROM user_reactions ur
WHERE
    ur.message_id = $1
GROUP BY
    ur.reaction_type;

-- Room Deletion Operations
-- name: DeleteRoomAndMessages :execrows
WITH
    room_check AS (
        SELECT 1
        FROM
            room_creators rc
            JOIN user_sessions us ON rc.creator_session_id = us.id
        WHERE
            rc.room_id = $1
            AND us.session_token = $2
    ),
    deleted_reactions AS (
        DELETE FROM user_reactions
        WHERE
            room_id = $1
            AND EXISTS (
                SELECT 1
                FROM room_check
            ) RETURNING user_reactions.id
    ),
    deleted_messages AS (
        DELETE FROM messages
        WHERE
            room_id = $1
            AND EXISTS (
                SELECT 1
                FROM room_check
            ) RETURNING messages.id
    )
DELETE FROM rooms
WHERE
    rooms.id = $1
    AND EXISTS (
        SELECT 1
        FROM room_check
    );