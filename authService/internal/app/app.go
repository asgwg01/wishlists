package app

import (
	redisImpl "authService/internal/cache/redis"
	"authService/internal/config"
	"authService/internal/grpc/handler"
	jwtservice "authService/internal/services/jwtService"
	userservice "authService/internal/services/userService"
	"authService/internal/storage/postgres"
	"errors"
	"log/slog"
	"net"
	"os"

	authv1 "github.com/asgwg01/wishlists/pkg/proto/auth/v1"

	"google.golang.org/grpc"
)

type App struct {
	log    *slog.Logger
	server *grpc.Server
	cfg    *config.Config
}

func New(log *slog.Logger, cfg *config.Config) *App {

	application := App{}
	application.log = log
	application.cfg = cfg

	// init storage
	storage, err := postgres.NewStorage(log, cfg.StorageConfig)
	if err != nil {
		log.Error("can't create storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	//storage, _ := stubstorage.NewStorage(log) // STUB storage

	//init redis
	cache, err := redisImpl.NewCache(log, &cfg.RedisConfig)
	if err != nil {
		log.Error("can't create cache", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// init services
	jwtService := jwtservice.New(log, &cfg.JWTConfig, cache)
	userService := userservice.New(log, storage, jwtService)

	// grpc
	application.server = grpc.NewServer()
	// init grpc handler
	grpcHandler := handler.NewGrpcHandler(*userService)
	authv1.RegisterAuthServiceServer(application.server, grpcHandler)

	return &application
}

func (a *App) Start() {
	const logPrefix = "app.Start"
	log := a.log.With(
		slog.String("where", logPrefix),
	)

	// Запускаем gRPC сервер
	lis, err := net.Listen("tcp", ":"+a.cfg.GRPCConfig.Port)
	if err != nil {
		log.Error("Failed to listen", slog.String("error", err.Error()))
	}

	log.Info("grpc server is runing", slog.String("addr", lis.Addr().String()))

	if err := a.server.Serve(lis); err != nil {
		if errors.Is(err, grpc.ErrServerStopped) {
			log.Info("grpc server is stoped")
		} else {
			log.Error("grpc server error", slog.String("error", err.Error()))
		}
	}
}

func (a *App) Stop() {
	const logPrefix = "app.Stop"
	log := a.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Server stoping")
	a.server.GracefulStop()
	log.Info("Stop grpc server end")

}
