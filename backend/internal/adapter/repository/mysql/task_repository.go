package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/a-s/connect-task-manage/internal/adapter/repository"
	"github.com/a-s/connect-task-manage/internal/domain/model"
	"github.com/a-s/connect-task-manage/internal/infrastructure/config"
	"github.com/a-s/connect-task-manage/sql/query"
)

type taskRepository struct {
	db      *sql.DB
	queries *query.Queries
	tx      *sql.Tx // トランザクションを保持するフィールド
}

// NewTaskRepository は新しい TaskRepository の実装を返します。
func NewTaskRepository(cfg *config.Config) (repository.TaskRepository, error) { // *config.Config を受け取る
	db, err := sql.Open("mysql", cfg.DB.DSN) // cfg から DSN を取得
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &taskRepository{
		db:      db,
		queries: query.New(db),
		tx:      nil,
	}, nil
}

// トランザクションを開始
func (r *taskRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// トランザクション内での操作用 (Queries オブジェクトを返す)
func (r *taskRepository) WithTx(tx *sql.Tx) repository.TaskRepository {
	return &taskRepository{
		db:      r.db,                 // db は共通
		queries: r.queries.WithTx(tx), // WithTx で新しい Queries を作成
		tx:      tx,                   // tx を保持
	}
}

func (r *taskRepository) CreateTask(ctx context.Context, task *model.Task) error {
	return r.queries.CreateTask(ctx, &query.CreateTaskParams{
		ID:          task.ID,
		Title:       task.Title,
		Description: sql.NullString{String: task.Description, Valid: task.Description != ""},
		IsCompleted: task.IsCompleted,
		UserID:      task.UserID,
	})
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	err := r.queries.UpdateTask(ctx, &query.UpdateTaskParams{
		ID:          task.ID,
		Title:       task.Title,
		Description: sql.NullString{String: task.Description, Valid: task.Description != ""},
		IsCompleted: task.IsCompleted,
	})
	if err != nil {
		return nil, err
	}
	updatedTask, err := r.queries.GetTaskByID(ctx, task.ID)
	if err != nil {
		return nil, err
	}

	return &model.Task{ // domain modelに変換
		ID:          updatedTask.ID,
		Title:       updatedTask.Title,
		Description: updatedTask.Description.String, // Stringを取り出す
		IsCompleted: updatedTask.IsCompleted,
		UserID:      updatedTask.UserID,
		CreatedAt:   updatedTask.CreatedAt,
		UpdatedAt:   updatedTask.UpdatedAt,
	}, nil
}

func (r *taskRepository) ListTasks(ctx context.Context, userID string) ([]*model.Task, error) {
	queryTasks, err := r.queries.ListTasks(ctx, userID)
	if err != nil {
		return nil, err
	}

	tasks := []*model.Task{}
	for _, t := range queryTasks { // queryTasks を range でループ
		tasks = append(tasks, &model.Task{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description.String,
			IsCompleted: t.IsCompleted,
			UserID:      t.UserID,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}
	return tasks, nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, id string) error {
	return r.queries.DeleteTask(ctx, id)
}
func (r *taskRepository) GetTaskByID(ctx context.Context, id string) (*model.Task, error) {
	task, err := r.queries.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}
	return &model.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description.String, // Stringを取り出す
		IsCompleted: task.IsCompleted,
		UserID:      task.UserID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}
