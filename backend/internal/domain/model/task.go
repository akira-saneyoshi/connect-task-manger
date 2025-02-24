package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Priority はタスクの優先度を表す型
type Priority string

// 優先度の定数
const (
	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityLow    Priority = "low"
)

type Task struct {
	ID          string
	Title       string
	Description string
	IsCompleted bool
	UserID      string  // Taskの作成者
	AssigneeID  *string // Taskの担当者
	Priority    Priority
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTask は新しい User エンティティを作成します。
func NewTask(title, description string, userID string, priority Priority, dueDate *time.Time) (*Task, error) {

	//priorityのバリデーション
	switch priority {
	case PriorityHigh, PriorityMedium, PriorityLow:
	default:
		return nil, fmt.Errorf("invalid priority: %v", priority)
	}

	return &Task{
		ID:          uuid.NewString(),
		Title:       title,
		Description: description,
		IsCompleted: false,
		UserID:      userID,
		AssigneeID:  nil,
		Priority:    priority,
		DueDate:     dueDate,
	}, nil
}

// Update はタスクの情報を更新します。
func (t *Task) Update(title, description string, isCompleted bool, assigneeID *string, priority Priority, dueDate *time.Time) error {
	//priorityのバリデーション
	switch priority {
	case PriorityHigh, PriorityMedium, PriorityLow:
	default:
		return fmt.Errorf("invalid priority: %v", priority)
	}

	t.Title = title
	t.Description = description
	t.IsCompleted = isCompleted
	t.AssigneeID = assigneeID
	t.Priority = priority
	t.DueDate = dueDate
	return nil
}
