-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2
)
RETURNING *;

-- name: GetUser :one
SELECT id, email, created_at, updated_at
FROM Users
WHERE id = $1;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT id, email, created_at, updated_at, hashed_password
FROM Users
WHERE email = $1;