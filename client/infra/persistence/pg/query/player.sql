-- name: PlayerByID :one
SELECT * FROM players
WHERE id = $1 LIMIT 1;

-- name: PlayerByEmailOrUsername :one
SELECT * FROM players
WHERE email = sqlc.arg(email_or_username) OR
    username = sqlc.arg(email_or_username) LIMIT 1;

-- name: CreatePlayer :one
INSERT INTO players (
    username, email, password
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: VerifyEmail :exec
UPDATE players SET email_verified_at = NOW()
WHERE id = $1;

-- name: UpdatePassword :exec
UPDATE players SET password = $1, password_updated_at = NOW()
WHERE id = $2;
