-- name: CreateUser :one
INSERT INTO
    users (name, email, password)
VALUES ($1, $2, $3)
RETURNING
    *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdatePassword :execrows
UPDATE users SET password = $2, updated_at = now() WHERE id = $1;