package app

import (
	"context"
	"errors"
	"log/slog"
	"notificationService/internal/config"
	"notificationService/internal/consumer"
	email "notificationService/internal/notifier"
	"time"
)

type App struct {
	log           *slog.Logger
	kafkaConsumer *consumer.KafkaConsumer
	ctx           context.Context
	cancelCalback context.CancelFunc
}

func New(log *slog.Logger, cfg *config.Config) *App {

	emailNotifier := email.NewEmailNotifier(log, cfg)

	kafkaConsumer := consumer.NewKafkaConsumer(log, cfg.KafkaConfig, emailNotifier)

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		log:           log,
		kafkaConsumer: kafkaConsumer,
		ctx:           ctx,
		cancelCalback: cancel,
	}
}

func (a *App) Start() {
	const logPrefix = "app.Start"
	log := a.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("start server")

	if err := a.kafkaConsumer.Start(a.ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Info("server is stoped")
		} else {
			log.Error("server error", slog.String("error", err.Error()))
		}
	}
}

func (a *App) Stop() {
	const logPrefix = "app.Stop"
	log := a.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("server stoping")

	a.cancelCalback()
	time.Sleep(2 * time.Second)

	a.kafkaConsumer.Close()

}
