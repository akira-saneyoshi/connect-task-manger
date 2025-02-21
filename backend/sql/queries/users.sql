-- sql/queries/users.sql

-- name: CreateUser :exec
INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?);

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: UpdateUser :exec
UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?;