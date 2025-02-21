package repository

import (
	"context"

	"github.com/a-s/connect-task-manage/internal/domain/model"
)

// UserRepository はユーザーデータへのアクセスを抽象化するインターフェースです。
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
}
