-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password
FROM users
WHERE email = $1;

-- name: UpdateUserCredentials :exec
UPDATE users
SET 
    email = $1,
    hashed_password = $2,
    updated_at = Now()
WHERE id = $3;

-- name: GetRoleByID :one
SELECT role FROM users WHERE id = $1;