package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// init config
	//config := config.LoadConfig()

	// init loger
	log := SetupLoger()
	log.Info("starting logger")

	// init app
	//application := app.New(log, *config, storage, service)
	//go application.Start()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signal := <-stop

	//application.Stop()

	log.Info("app stopped", slog.String("reason signal", signal.String()))
}

func SetupLoger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}
