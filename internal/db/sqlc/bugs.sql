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


-- name: GetAllBugs :many
SELECT * FROM bugs
ORDER BY created_at DESC;

-- name: GetBugsByID :one
SELECT * FROM bugs
WHERE Id = $1;


-- name: UpdateBugByID :exec
UPDATE bugs
SET 
    title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    updated_at = Now()
WHERE id = $1;

-- name: DeleteBugByID :exec
DELETE FROM bugs
WHERE id = $1;

