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
	"github.com/a-s/connect-task-manage/internal/infrastructure/logger"
	"github.com/a-s/connect-task-manage/pkg/authorization"
	"github.com/a-s/connect-task-manage/pkg/logging"
	"go.uber.org/fx"
	"go.uber.org/zap"
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

// NewUserServiceServer は UserServiceServer のコンストラクタ (Fx 用)
func NewUserServiceServer(userService *service.UserService) *UserServiceServer {
	return &UserServiceServer{userService: userService}
}

// InterceptorFuncToInterceptor は connect.UnaryInterceptorFunc を connect.Interceptor に変換
func InterceptorFuncToInterceptor(fn connect.UnaryInterceptorFunc) connect.Interceptor {
	return connect.Interceptor(fn)
}

// NewInterceptors はインターセプターのリストを提供
func NewInterceptors(interceptors []connect.UnaryInterceptorFunc) []connect.Interceptor {
	converted := make([]connect.Interceptor, len(interceptors))
	for i, interceptor := range interceptors {
		converted[i] = connect.Interceptor(interceptor)
	}
	return converted
}

// NewHTTPServer は HTTP サーバーのコンストラクタ
func NewHTTPServer(
	lc fx.Lifecycle,
	cfg *config.Config,
	userServiceServer *UserServiceServer,
	log *zap.Logger,
	interceptors []connect.Interceptor,
) *http.Server {
	services := []string{userv1connect.UserServiceName}
	reflector := grpcreflect.NewStaticReflector(services...)

	mux := http.NewServeMux()
	path, handler := userv1connect.NewUserServiceHandler(
		userServiceServer,
		connect.WithInterceptors(interceptors...),
	)
	mux.Handle(path, handler)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting server", zap.String("port", cfg.App.Port))
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal("Server listen error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down server...")
			return server.Shutdown(ctx)
		},
	})

	return server
}

func main() {
	app := fx.New(
		fx.Provide(
			config.LoadConfig,
			logger.NewLogger,
			mysql.NewUserRepository,
			jwt.NewJWTManager,
			service.NewUserService,
			NewUserServiceServer,
			fx.Annotate(
				authorization.NewAuthInterceptor,
				fx.ResultTags(`group:"interceptors"`),
			),
			fx.Annotate(
				logging.NewLoggingInterceptor,
				fx.ResultTags(`group:"interceptors"`),
			),
			fx.Annotate(
				NewInterceptors,
				fx.ParamTags(`group:"interceptors"`),
			),
			NewHTTPServer,
		),
		fx.Invoke(func(server *http.Server) {}),
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Start(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		log.Fatal(err)
	}
}
