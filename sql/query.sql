-- name: GetTask :one
SELECT * FROM tasks WHERE id = $1;

-- name: ListTasks :many
SELECT * FROM tasks ORDER BY id;

-- name: CreateTask :one
INSERT INTO tasks (title, completed)
VALUES ($1, $2)
    RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;