package main

import (
	"fmt"
	"io"
	"log/slog"
	"notificationService/internal/app"
	"notificationService/internal/config"

	stdLog "log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asgwg01/wishlists/pkg/logger"
	"github.com/logbull/logbull-go/logbull"
)

func main() {
	// init config
	config := config.LoadConfig()

	// init loger
	log, logFile := SetupLoger(config)
	defer logFile.Close()
	log.Info("starting logger")

	// // init app
	application := app.New(log, config)

	go application.Start()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signal := <-stop

	application.Stop()

	log.Info("app stopped", slog.String("reason signal", signal.String()))
}

const (
	envLocal       = "local"
	envDev         = "dev"
	envProd        = "prod"
	serviceName    = "notificationService"
	logFileName    = "notificationService.log"
	logFileBigSize = 1024 * 1024
)

func SetupLoger(cfg *config.Config) (*slog.Logger, *os.File) {
	var log *slog.Logger

	lfInfo, err := os.Stat(logFileName)
	if err != nil {
		if !os.IsNotExist(err) {
			stdLog.Fatalf("can not create or open %s: %s", logFileName, err.Error())
		}
	}

	if !os.IsNotExist(err) && lfInfo.Size() > logFileBigSize {
		newName := "before_" + time.Now().Local().Format(time.DateTime) + "_" + logFileName
		err = os.Rename(logFileName, newName)
		if err != nil {
			stdLog.Fatalf("can not rename %s", logFileName)
		}
	}

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		stdLog.Fatalf("can not create or open %s: %s", logFileName, err.Error())
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	var handler slog.Handler

	switch cfg.Env {
	case envProd:
		handler = slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelInfo})
	case envDev:
		handler = slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug})
	case envLocal:
		fallthrough
	default:
		handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug})
	}

	if cfg.LogBullConfig.URL != "" {
		bullHandler, err := logbull.NewSlogHandler(logbull.Config{
			Host:      "http://" + cfg.LogBullConfig.URL + ":" + cfg.LogBullConfig.Port,
			ProjectID: cfg.LogBullConfig.ProjectID,
			APIKey:    "",
		})
		if err != nil {
			panic(err)
		}

		// // Объединяем handlers - логи идут и в файл/консоль, и в Log Bull
		handler = &logger.MultiHandler{
			Handlers: []slog.Handler{handler, bullHandler},
		}
	}

	log = slog.New(handler)

	log = log.With(slog.String("service_name", serviceName))

	fmt.Fprintln(logFile, "\n\n\nNew Session\n=====================================================")

	return log, logFile
}
