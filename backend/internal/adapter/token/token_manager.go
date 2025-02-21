package token

import "github.com/a-s/connect-task-manage/internal/domain/model"

// TokenManager はトークンの生成と検証を抽象化するインターフェースです。
type TokenManager interface {
	Generate(user *model.User) (string, error)
	Verify(token string) (string, error) // (userID, error) を返す
}
