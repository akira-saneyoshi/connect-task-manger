package logging

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/a-s/connect-task-manage/internal/infrastructure/logger" // loggerパッケージをimport
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// NewLoggingInterceptor はロギングインターセプターを作成します。
func NewLoggingInterceptor(log *zap.Logger) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			// リクエストIDの生成 (存在しない場合)
			requestID := req.Header().Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// リクエストID付きのロガーをコンテキストに追加
			ctx = context.WithValue(ctx, logger.LoggerKey, logger.WithRequestID(log, requestID)) // LoggerKey を使用
			log = logger.WithRequestID(log, requestID)                                           //リクエストID付きのロガーを作成

			start := time.Now()

			// リクエスト情報のログ出力
			log.Info("request started",
				zap.String("method", req.Spec().Procedure),
				zap.String("request_id", requestID),
				zap.Any("headers", req.Header()), // ヘッダーもログ出力
			)

			// リクエストの実行
			res, err := next(ctx, req)

			duration := time.Since(start)

			// レスポンス情報のログ出力 (エラーの有無でレベルを分ける)
			if err != nil {
				log.Error("request failed",
					zap.String("method", req.Spec().Procedure),
					zap.String("request_id", requestID),
					zap.Duration("duration", duration),
					zap.Error(err),
				)
			} else {
				log.Info("request completed",
					zap.String("method", req.Spec().Procedure),
					zap.String("request_id", requestID),
					zap.Duration("duration", duration),
					zap.Any("headers", res.Header()), //レスポンスヘッダー
				)
			}

			return res, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
