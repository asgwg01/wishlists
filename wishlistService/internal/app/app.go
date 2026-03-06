package app

import (
	"errors"
	"log/slog"
	"net"
	"os"
	"wishlistService/internal/config"
	"wishlistService/internal/grpc/client"
	"wishlistService/internal/grpc/handler"
	"wishlistService/internal/producer"
	itemservice "wishlistService/internal/services/itemService"
	"wishlistService/internal/services/wishlistService"
	"wishlistService/internal/storage/postgres"

	wishlistv1 "github.com/asgwg01/wishlists/pkg/proto/wishlists/v1"

	"google.golang.org/grpc"
)

type App struct {
	log           *slog.Logger
	server        *grpc.Server
	authClient    *client.AuthClient
	kafkaPronucer *producer.KafkaProducer
	cfg           *config.Config
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

	// grpc client
	authClient, err := client.NewAuthClient(cfg.GRPCConfig.AuthServiceAddr, cfg.GRPCConfig.AuthServicePort)
	if err != nil {
		log.Error("Failed to connect to auth service", slog.String("error", err.Error()))
		os.Exit(1)
	}
	application.authClient = authClient

	// Создаем Kafka producer
	kafkaProducer := producer.NewKafkaProducer(
		log,
		cfg.KafkaConfig.BrokerUrl+":"+cfg.KafkaConfig.BrokerPort,
		cfg.KafkaConfig.Topic,
	)
	application.kafkaPronucer = kafkaProducer

	// init services
	wishlistService := wishlistService.New(log, storage)
	itemService := itemservice.New(log, storage, storage, authClient, kafkaProducer)

	// grpc
	application.server = grpc.NewServer()
	// init grpc handler
	grpcHandler := handler.NewGrpcHandler(wishlistService, itemService)
	wishlistv1.RegisterWishlistServiceServer(application.server, grpcHandler)

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
	a.authClient.Close()
	a.kafkaPronucer.Close()
	a.server.GracefulStop()
	log.Info("Stop server end")

}
