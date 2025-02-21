package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/a-s/connect-task-manage/internal/adapter/repository"
	"github.com/a-s/connect-task-manage/internal/domain/model"
	"github.com/a-s/connect-task-manage/internal/infrastructure/config"
	"github.com/a-s/connect-task-manage/sql/query"

	_ "github.com/go-sql-driver/mysql"
)

type userRepository struct {
	db      *sql.DB
	queries *query.Queries
}

// NewUserRepository は新しい UserRepository の実装を返します
func NewUserRepository(cfg *config.Config) (repository.UserRepository, error) {
	db, err := sql.Open("mysql", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &userRepository{
		db:      db,
		queries: query.New(), // ここを修正：引数なし
	}, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.queries.CreateUser(ctx, r.db, &query.CreateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, r.db, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := r.queries.GetUserByID(ctx, r.db, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}
	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.queries.UpdateUser(ctx, r.db, &query.UpdateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})

	if err != nil {
		return nil, err
	}

	updatedUser, err := r.queries.GetUserByID(ctx, r.db, user.ID)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:        updatedUser.ID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		Password:  updatedUser.Password,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}
