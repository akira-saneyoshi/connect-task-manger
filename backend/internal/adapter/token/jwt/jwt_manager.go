package jwt

import (
	"fmt"
	"time"

	"github.com/a-s/connect-task-manage/internal/adapter/token"
	"github.com/a-s/connect-task-manage/internal/domain/model"
	"github.com/a-s/connect-task-manage/internal/infrastructure/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// NewJWTManager は新しい JWTManager インスタンスを作成します。
func NewJWTManager(cfg *config.Config) token.TokenManager {
	return &JWTManager{
		secretKey:     cfg.JWT.Secret,
		tokenDuration: time.Duration(cfg.JWT.DurationMinutes) * time.Minute,
	}
}

// Claims は JWT のカスタムクレームを表します。
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// Generate はユーザー情報に基づいて JWT トークンを生成します。
func (m *JWTManager) Generate(user *model.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ddd-auth-app",
			Subject:   user.Email, // サブジェクトは email とする
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// Verify は JWT トークンを検証し、ユーザー ID を返します。
func (m *JWTManager) Verify(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}
	return "", model.ErrInvalidToken
}
