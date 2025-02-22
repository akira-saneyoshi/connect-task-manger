-- sql/queries/tasks.sql

-- name: CreateTask :exec
INSERT INTO tasks (id, title, description, is_completed, user_id)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateTask :exec
UPDATE tasks SET title = ?, description = ?, is_completed = ? WHERE id = ?;

-- name: ListTasks :many
SELECT * FROM tasks WHERE user_id = ? ORDER BY created_at DESC;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = ?;

-- name: GetTaskByID :one
SELECT * FROM tasks WHERE id = ? LIMIT 1;