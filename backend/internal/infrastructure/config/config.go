package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config はアプリケーション全体の設定を保持します。
type Config struct {
	DB  DBConfig
	JWT JWTConfig
	App AppConfig
}

// DBConfig はデータベース接続設定を保持します。
type DBConfig struct {
	DSN string
}

// JWTConfig は JWT 認証関連の設定を保持します。
type JWTConfig struct {
	Secret          string
	DurationMinutes int
}

// AppConfig はアプリケーションの基本設定を保持します。
type AppConfig struct {
	Port string
}

// LoadConfig は .env ファイルおよび環境変数から設定を読み込みます。
func LoadConfig() (*Config, error) {
	// .env ファイルを読み込む (存在する場合)
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Println("Error loading .env file, using environment variables")
	}

	// 環境変数から設定値を読み込む
	dbDSN := getEnv("DB_DSN", "user:password@tcp(localhost:3306)/authdb?parseTime=true")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	jwtDurationMinutes, err := getEnvInt("JWT_DURATION_MINUTES", 15)
	if err != nil {
		return nil, err
	}
	appPort := getEnv("APP_PORT", "8080")

	return &Config{
		DB: DBConfig{
			DSN: dbDSN,
		},
		JWT: JWTConfig{
			Secret:          jwtSecret,
			DurationMinutes: jwtDurationMinutes,
		},
		App: AppConfig{
			Port: appPort,
		},
	}, nil
}

// getEnv は指定された環境変数の値を取得します。
// 環境変数が設定されていない場合は、デフォルト値を返します。
func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// getEnvInt は指定された環境変数の値を整数として取得します。
func getEnvInt(key string, defaultValue int) (int, error) {
	valueStr := getEnv(key, fmt.Sprintf("%d", defaultValue))
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid integer value for %s: %w", key, err)
	}
	return value, nil

}
