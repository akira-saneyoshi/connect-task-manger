package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	userv1 "github.com/a-s/connect-task-manage/gen/api/user/v1"
	"github.com/a-s/connect-task-manage/gen/api/user/v1/userv1connect"
	"github.com/a-s/connect-task-manage/internal/adapter/repository/mysql"
	"github.com/a-s/connect-task-manage/internal/adapter/token/jwt"
	"github.com/a-s/connect-task-manage/internal/domain/service"
	"github.com/a-s/connect-task-manage/internal/infrastructure/config"
	"github.com/a-s/connect-task-manage/pkg/authorization"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type UserServiceServer struct {
	userService *service.UserService
}

// CreateUser handles the CreateUser RPC method.
func (s *UserServiceServer) CreateUser(
	ctx context.Context,
	req *connect.Request[userv1.CreateUserRequest],
) (*connect.Response[userv1.CreateUserResponse], error) {

	_, err := s.userService.CreateUser(ctx, req.Msg.Name, req.Msg.Email, req.Msg.Password)
	if err != nil {
		// エラーの種類に応じて適切なエラーコードを返す
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := connect.NewResponse(&userv1.CreateUserResponse{
		User: nil, // User情報を含めない
	})
	return res, nil
}

func (s *UserServiceServer) Login(
	ctx context.Context,
	req *connect.Request[userv1.LoginRequest],
) (*connect.Response[userv1.LoginResponse], error) {
	token, err := s.userService.Login(ctx, req.Msg.Email, req.Msg.Password)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}
	res := connect.NewResponse(&userv1.LoginResponse{
		AccessToken: token,
	})

	return res, nil
}

func (s *UserServiceServer) UpdateUser(
	ctx context.Context,
	req *connect.Request[userv1.UpdateUserRequest],
) (*connect.Response[userv1.UpdateUserResponse], error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("user id not found in context"))
	}

	updatedUser, err := s.userService.UpdateUser(ctx, userID, req.Msg.Name, req.Msg.Email, req.Msg.Password)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&userv1.UpdateUserResponse{
		User: &userv1.User{
			Id:        updatedUser.ID,
			Name:      updatedUser.Name,
			Email:     updatedUser.Email,
			CreatedAt: updatedUser.CreatedAt.Format(time.RFC3339),
			UpdatedAt: updatedUser.UpdatedAt.Format(time.RFC3339),
		},
	})

	return res, nil
}

func (s *UserServiceServer) Logout(
	ctx context.Context,
	req *connect.Request[userv1.LogoutRequest],
) (*connect.Response[userv1.LogoutResponse], error) {

	res := connect.NewResponse(&userv1.LogoutResponse{})
	return res, nil

}

func (s *UserServiceServer) GetMe(
	ctx context.Context,
	req *connect.Request[userv1.GetMeRequest],
) (*connect.Response[userv1.GetMeResponse], error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("user id not found in context"))
	}

	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&userv1.GetMeResponse{
		User: &userv1.User{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	})

	return res, nil
}

func main() {
	// 設定の読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	//依存性の注入
	userRepository, err := mysql.NewUserRepository(cfg)
	if err != nil {
		log.Fatal(err)
	}

	tokenManager := jwt.NewJWTManager(cfg)

	userService := service.NewUserService(userRepository, tokenManager)

	// サーバーの準備
	userServiceServer := &UserServiceServer{
		userService: userService,
	}

	// インターセプターの作成
	authInterceptor := authorization.NewAuthInterceptor(tokenManager)

	// connect-go のハンドラを生成
	mux := http.NewServeMux()
	path, handler := userv1connect.NewUserServiceHandler(
		userServiceServer,
		connect.WithInterceptors(authInterceptor), //インターセプター追加
	)
	mux.Handle(path, handler)

	// サーバーの起動 (graceful shutdown 付き)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: h2c.NewHandler(mux, &http2.Server{}), // h2c を有効化
	}

	go func() {
		fmt.Printf("... Listening on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// シグナルを待機して graceful shutdown を行う
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	fmt.Println("Server exiting")

}
