package service

import (
	"context"
	"database/sql"
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

// WithTx はトランザクションを開始し、トランザクション内で操作を行うための新しい UserService インスタンスを返します。
func (s *UserService) WithTx(tx *sql.Tx) *UserService {
	return &UserService{
		userRepository: s.userRepository.WithTx(tx), // トランザクション用のリポジトリを使用
		tokenManager:   s.tokenManager,              // tokenManager は共通
	}
}

func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*model.User, error) {
	existingUser, _ := s.userRepository.GetUserByEmail(ctx, email)
	if existingUser != nil {
		return nil, model.ErrUserAlreadyExists
	}

	id := uuid.New().String()
	user, err := model.NewUser(id, name, email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user entity: %w", err)
	}

	_, err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in repository: %w", err)
	}
	return nil, nil
}

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

func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser はトランザクション内でユーザー情報を更新します。
func (s *UserService) UpdateUser(ctx context.Context, id, name, email, password string) (*model.User, error) {

	// トランザクション開始 (WithTx が呼ばれていない場合)
	tx, err := s.userRepository.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	// トランザクション用の UserService を作成
	txService := s.WithTx(tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // 再度パニックさせる
		} else if err != nil {
			_ = tx.Rollback() // エラーが発生したらロールバック
		} else {
			err = tx.Commit() // 成功したらコミット
		}
	}()

	user, err := txService.userRepository.GetUserByID(ctx, id) //txServiceを使う
	if err != nil {
		return nil, err
	}

	if err := user.Update(name, email, password); err != nil {
		return nil, err
	}

	updatedUser, err := txService.userRepository.UpdateUser(ctx, user) //txServiceを使う
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
