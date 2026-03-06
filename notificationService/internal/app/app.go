package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"notificationService/internal/config"
	"notificationService/internal/http/handlers/info"
	dowork "notificationService/internal/services/doWork"
	"notificationService/internal/storage"

	"github.com/gorilla/mux"
)

type App struct {
	log     *slog.Logger
	server  *http.Server
	storage storage.IStorage
}

func New(log *slog.Logger, cfg config.Config, storage storage.IStorage, service dowork.IServiceWorkSome) *App {

	router := mux.NewRouter()
	router.HandleFunc("/info", info.NewHandler(log, service)).Methods("GET")

	server := &http.Server{
		Addr:         cfg.Addres,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &App{
		log:     log,
		server:  server,
		storage: storage,
	}
}

func (a *App) Start() {
	const logPrefix = "app.Start"
	log := a.log.With(
		slog.String("where", logPrefix),
		slog.String("host", a.server.Addr),
	)

	log.Info("start server")
	if err := a.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
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

	if err := a.server.Shutdown(context.Background()); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			// log.Info("server is stoped")
			// will be printed in app.Start
		} else {
			log.Error("server error", slog.String("error", err.Error()))
		}
	}

}
