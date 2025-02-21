package authorization

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/a-s/connect-task-manage/internal/adapter/token"
	"github.com/a-s/connect-task-manage/internal/domain/model"
)

// NewAuthInterceptor は認証インターセプターを作成します。
func NewAuthInterceptor(tm token.TokenManager) connect.UnaryInterceptorFunc { // 戻り値の型は connect.UnaryInterceptorFunc
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			// 認証が不要なメソッドのリスト (Login, CreateUser など)
			unprotectedMethods := map[string]struct{}{
				"/user.v1.UserService/CreateUser": {},
				"/user.v1.UserService/Login":      {},
			}

			// メソッド名を取得
			methodName := req.Spec().Procedure
			if _, ok := unprotectedMethods[methodName]; ok {
				// 認証不要なメソッドはそのまま next に渡す
				return next(ctx, req)
			}

			// Authorization ヘッダーからトークンを取得
			authHeader := req.Header().Get("Authorization")
			if authHeader == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, model.ErrUnauthorized)
			}

			// "Bearer " プレフィックスを取り除く
			tokenString := ""
			_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, model.ErrUnauthorized)
			}

			// トークンを検証し、ユーザー ID を取得
			userID, err := tm.Verify(tokenString)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("token verification failed: %w", err))
			}

			// コンテキストにユーザー ID を設定
			ctx = context.WithValue(ctx, "userID", userID)

			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor) // connect.UnaryInterceptorFunc 型の値を返す
}
