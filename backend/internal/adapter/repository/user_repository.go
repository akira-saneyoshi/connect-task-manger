package repository

import (
	"context"
	"database/sql"

	"github.com/a-s/connect-task-manage/internal/domain/model"
	// sqlc の生成コード
)

// UserRepository はユーザーデータへのアクセスを抽象化するインターフェースです。
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)

	// トランザクション関連のメソッド
	BeginTx(ctx context.Context) (*sql.Tx, error) // トランザクション開始
	WithTx(tx *sql.Tx) UserRepository             // トランザクション内での操作用
}
