// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package query

import (
	"database/sql"
	"time"
)

type Task struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	IsCompleted bool           `json:"is_completed"`
	UserID      string         `json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
