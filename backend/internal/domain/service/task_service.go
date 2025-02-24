package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/a-s/connect-task-manage/internal/adapter/repository"
	"github.com/a-s/connect-task-manage/internal/domain/model"
)

type TaskService struct {
	taskRepository repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) *TaskService {
	return &TaskService{taskRepository: taskRepo}
}

func (s *TaskService) WithTx(tx *sql.Tx) *TaskService {
	return &TaskService{
		taskRepository: s.taskRepository.WithTx(tx),
	}
}

func (s *TaskService) CreateTask(ctx context.Context, title, description, userID string, priority string, dueDate *time.Time) error {
	task, err := model.NewTask(title, description, userID, model.Priority(priority), dueDate) // model.Priority に変換
	if err != nil {
		return err
	}
	return s.taskRepository.CreateTask(ctx, task)
}

func (s *TaskService) UpdateTask(ctx context.Context, id, title, description string, isCompleted bool, assigneeID *string, priority string, dueDate *time.Time) (*model.Task, error) {

	task, err := s.taskRepository.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := task.Update(title, description, isCompleted, assigneeID, model.Priority(priority), dueDate); err != nil { // model.Priorityに変換
		return nil, err
	}
	return s.taskRepository.UpdateTask(ctx, task)
}

func (s *TaskService) ListTasks(ctx context.Context, userID string) ([]*model.Task, error) {
	return s.taskRepository.ListTasks(ctx, userID)
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	return s.taskRepository.DeleteTask(ctx, id)
}
func (s *TaskService) GetTaskByID(ctx context.Context, id string) (*model.Task, error) {
	task, err := s.taskRepository.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrUserNotFound //ドメイン層のエラーを返す
	}

	return task, nil
}
