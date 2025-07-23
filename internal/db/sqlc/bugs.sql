-- name: CreateBug :one
INSERT INTO bugs (id, title, description, posted_by, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    NOW(),
    NOW()
    
)
RETURNING *;