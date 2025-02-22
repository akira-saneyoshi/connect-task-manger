package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          string
	Title       string
	Description string
	IsCompleted bool
	UserID      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTask(title, description string, userID string) (*Task, error) {
	return &Task{
		ID:          uuid.NewString(),
		Title:       title,
		Description: description,
		IsCompleted: false,
		UserID:      userID,
	}, nil
}

func (t *Task) Update(title, description string, isCompleted bool) {
	t.Title = title
	t.Description = description
	t.IsCompleted = isCompleted
}
