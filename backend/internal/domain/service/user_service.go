package service

import (
	"context"
	"fmt"

	"github.com/a-s/connect-task-manage/internal/adapter/repository"
	"github.com/a-s/connect-task-manage/internal/adapter/token"
	"github.com/a-s/connect-task-manage/internal/domain/model"
	"github.com/google/uuid"
)

// UserService はユーザーに関するビジネスロジックを提供します。
type UserService struct {
	userRepository repository.UserRepository
	tokenManager   token.TokenManager
}

// NewUserService は新しい UserService インスタンスを作成します。
func NewUserService(userRepo repository.UserRepository, tokenManager token.TokenManager) *UserService {
	return &UserService{
		userRepository: userRepo,
		tokenManager:   tokenManager,
	}
}

// CreateUser は新しいユーザーを作成します。
func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*model.User, error) {
	// 既存ユーザーのチェック (email の一意性)
	existingUser, _ := s.userRepository.GetUserByEmail(ctx, email) //エラーを無視して存在チェック
	if existingUser != nil {
		return nil, model.ErrUserAlreadyExists
	}

	id := uuid.New().String()
	user, err := model.NewUser(id, name, email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user entity: %w", err)
	}

	_, err = s.userRepository.CreateUser(ctx, user) // 戻り値は使わない
	if err != nil {
		return nil, fmt.Errorf("failed to create user in repository: %w", err)
	}
	return nil, nil // repository の戻り値が nil なので、ここも nil を返す。
}

// Login はユーザーを認証し、アクセストークンを生成します。
func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to get user by email: %w", err)
	}

	if err := user.Authenticate(password); err != nil {
		return "", model.ErrAuthentication
	}

	token, err := s.tokenManager.Generate(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// GetUserByID は、IDでユーザーを取得します。
func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser はユーザー情報を更新します。
func (s *UserService) UpdateUser(ctx context.Context, id, name, email, password string) (*model.User, error) {

	user, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := user.Update(name, email, password); err != nil {
		return nil, err
	}

	updatedUser, err := s.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
