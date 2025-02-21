package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger はアプリケーション用のロガーを初期化します。
func NewLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // ISO8601 形式
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)          // デバッグレベル以上を出力

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

// WithRequestID はリクエスト ID を持つ新しいロガーを作成します。
func WithRequestID(logger *zap.Logger, requestID string) *zap.Logger {
	return logger.With(zap.String("request_id", requestID))
}

// ctxKey はコンテキストで使用するキーの型です。
type ctxKey int

const (
	LoggerKey ctxKey = iota // 公開 (LoggerKey)
)
