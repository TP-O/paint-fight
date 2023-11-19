-- name: PlayerByID :one
SELECT * FROM players
WHERE user_id = $1;

-- name: PlayersByUsername :many
SELECT user_id, username FROM players
WHERE username LIKE sqlc.arg(username)::varchar || '%';

-- name: CreatePlayer :one
INSERT INTO players (
    user_id, username
) VALUES (
    $1, $2
)
RETURNING *;
