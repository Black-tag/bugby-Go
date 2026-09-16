-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, username, role_id, password_hash)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, password_hash
FROM users
WHERE email = $1;

-- name: UpdateUserCredentials :exec
UPDATE users
SET 
    email = $1,
    password_hash = $2,
    updated_at = Now()
WHERE id = $3;

-- name: GetRoleByID :one
SELECT role_id FROM users WHERE id = $1;


-- name: GetAllUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: DeleteuserByID :exec
DELETE FROM users
WHERE id = $1;
