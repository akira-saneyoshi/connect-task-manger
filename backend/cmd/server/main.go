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
	"connectrpc.com/grpcreflect"
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

// CreateUser, Login, UpdateUser, Logout, GetMe メソッドは変更なし (省略)
func (s *UserServiceServer) CreateUser(
	ctx context.Context,
	req *connect.Request[userv1.CreateUserRequest],
) (*connect.Response[userv1.CreateUserResponse], error) {

	_, err := s.userService.CreateUser(ctx, req.Msg.Name, req.Msg.Email, req.Msg.Password)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := connect.NewResponse(&userv1.CreateUserResponse{
		User: nil,
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

	userRepository, err := mysql.NewUserRepository(cfg)
	if err != nil {
		log.Fatal(err)
	}
	tokenManager := jwt.NewJWTManager(cfg)
	userService := service.NewUserService(userRepository, tokenManager)
	userServiceServer := &UserServiceServer{userService: userService}
	authInterceptor := authorization.NewAuthInterceptor(tokenManager)

	// リフレクション用のサービス名リストを作成
	services := []string{
		userv1connect.UserServiceName,
	}
	reflector := grpcreflect.NewStaticReflector(services...) // 変更

	mux := http.NewServeMux()
	// connect-go のハンドラを登録
	path, handler := userv1connect.NewUserServiceHandler(
		userServiceServer,
		connect.WithInterceptors(authInterceptor),
		// WithGRPC() は不要
	)
	mux.Handle(path, handler)
	// リフレクション用のハンドラを登録
	mux.Handle(grpcreflect.NewHandlerV1(reflector))      // 追加
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector)) // 追加

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	go func() {
		fmt.Println("... Listening on :", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

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
