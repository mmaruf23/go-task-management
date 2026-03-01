-- name: CreateTask :one
INSERT INTO
    tasks (user_id, title, description)
VALUES ($1, $2, $3)
RETURNING
    *;

-- name: ListTaskByUser :many
SELECT *
FROM tasks
WHERE
    user_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- name: CountTaskByUser :one
SELECT COUNT(*) FROM tasks WHERE user_id = $1;

-- name: GetTaskByID :one
SELECT * FROM tasks WHERE id = $1 AND user_id = $2;

-- name: UpdateTask :one
UPDATE tasks
SET
    title = $3,
    description = $4,
    updated_at = now()
WHERE
    id = $2
    and user_id = $1
RETURNING
    *;

-- name: UpdateStatus :execrows
UPDATE tasks SET status = $3 WHERE id = $2 AND user_id = $1;