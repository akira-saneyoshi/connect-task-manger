package repository

import (
	"context"
	"database/sql"

	"github.com/a-s/connect-task-manage/internal/domain/model"
	// sqlc の生成コード
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *model.Task) error
	UpdateTask(ctx context.Context, task *model.Task) (*model.Task, error)
	ListTasks(ctx context.Context, userID string) ([]*model.Task, error)
	DeleteTask(ctx context.Context, id string) error
	GetTaskByID(ctx context.Context, id string) (*model.Task, error)

	// トランザクション関連 (UserRepository からコピー)
	BeginTx(ctx context.Context) (*sql.Tx, error)
	WithTx(tx *sql.Tx) TaskRepository
}
