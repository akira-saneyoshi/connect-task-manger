// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package query

import (
	"context"
)

const createUser = `-- name: CreateUser :exec

INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)
`

type CreateUserParams struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// sql/queries/users.sql
func (q *Queries) CreateUser(ctx context.Context, arg *CreateUserParams) error {
	_, err := q.exec(ctx, q.createUserStmt, createUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Password,
	)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = ? LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := q.queryRow(ctx, q.getUserByEmailStmt, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = ? LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id string) (*User, error) {
	row := q.queryRow(ctx, q.getUserByIDStmt, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?
`

type UpdateUserParams struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	ID       string `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg *UpdateUserParams) error {
	_, err := q.exec(ctx, q.updateUserStmt, updateUser,
		arg.Name,
		arg.Email,
		arg.Password,
		arg.ID,
	)
	return err
}
