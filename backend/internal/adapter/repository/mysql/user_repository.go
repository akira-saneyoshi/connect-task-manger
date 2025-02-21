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
	tx      *sql.Tx // トランザクションを保持するフィールド
}

// NewUserRepository は新しい UserRepository の実装を返します。
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
		queries: query.New(db),
		tx:      nil,
	}, nil
}

// トランザクションを開始
func (r *userRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// トランザクション内での操作用 (Queries オブジェクトを返す)
func (r *userRepository) WithTx(tx *sql.Tx) repository.UserRepository {
	return &userRepository{
		db:      r.db,                 // db は共通
		queries: r.queries.WithTx(tx), // WithTx で新しい Queries を作成
		tx:      tx,                   // tx を保持
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.queries.CreateUser(ctx, &query.CreateUserParams{ // 戻り値は err のみ
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

	user, err := r.queries.GetUserByEmail(ctx, email)
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
	user, err := r.queries.GetUserByID(ctx, id)
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

	err := r.queries.UpdateUser(ctx, &query.UpdateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})

	if err != nil {
		return nil, err
	}

	updatedUser, err := r.queries.GetUserByID(ctx, user.ID)
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
