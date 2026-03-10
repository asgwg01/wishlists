package app

import (
	"context"
	"errors"
	"gateway/internal/config"
	"gateway/internal/grpc/client"
	"gateway/internal/http/handlers/auth"
	"gateway/internal/http/handlers/booking"
	"gateway/internal/http/handlers/item"
	"gateway/internal/http/handlers/wishlist"
	"gateway/internal/http/middleware"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	log    *slog.Logger
	server *http.Server
	//grpcServer        *grpc.Server
	authServiceClient     *client.AuthClient
	wishlistServiceClient *client.WishlistClient
}

func New(log *slog.Logger, cfg *config.Config) *App {

	application := App{}
	application.log = log

	// grpc clients
	authServiceClient, err := client.NewAuthClient(cfg.GRPCConfig.AuthServiceAddr, cfg.GRPCConfig.AuthServicePort)
	if err != nil {
		log.Error("Failed to connect to auth service", slog.String("error", err.Error()))
		os.Exit(1)
	}
	application.authServiceClient = authServiceClient

	wishlistServiceClient, err := client.NewWishlistClient(cfg.GRPCConfig.WishlistServiceAddr, cfg.GRPCConfig.WishlistServicePort)
	if err != nil {
		log.Error("Failed to connect to wishlist service", slog.String("error", err.Error()))
		os.Exit(1)
	}
	application.wishlistServiceClient = wishlistServiceClient

	// routers
	router := mux.NewRouter()

	// middleware
	rateLimitter := middleware.NewRateLimitter(cfg.RateLimitsConfig)

	router.Use(middleware.LoggingMiddleware(log))
	router.Use(middleware.CORSMiddleware(cfg.CORSConfig))
	router.Use(middleware.TraceMiddleware())
	router.Use(middleware.AuthMiddleware(log, authServiceClient))
	router.Use(middleware.RateLimitMiddleware(log, rateLimitter))

	//апи хендлеры
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	//аус
	authHandlers := auth.NewHandlers(log, authServiceClient)
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", authHandlers.Register).Methods("POST")
	authRouter.HandleFunc("/login", authHandlers.Login).Methods("POST")
	authRouter.HandleFunc("/self", authHandlers.GetCurrentUser).Methods("GET")
	authRouter.HandleFunc("/user/{user_id}", authHandlers.GetUserInfo).Methods("GET")

	//вишлисты
	wishlistHandlers := wishlist.NewHandlers(log, wishlistServiceClient)
	wishlistRouter := apiRouter.PathPrefix("/wishlists").Subrouter()
	wishlistRouter.HandleFunc("", wishlistHandlers.CreateWishlist).Methods("POST")
	wishlistRouter.HandleFunc("/public", wishlistHandlers.ListPublicWishlists).Methods("GET")
	wishlistRouter.HandleFunc("/{id}", wishlistHandlers.GetWishlist).Methods("GET")
	wishlistRouter.HandleFunc("/{id}", wishlistHandlers.UpdateWishlist).Methods("PUT")
	wishlistRouter.HandleFunc("/{id}", wishlistHandlers.DeleteWishlist).Methods("DELETE")
	wishlistRouter.HandleFunc("/user/{user_id}", wishlistHandlers.GetUserWishlists).Methods("GET")
	//айтемы в вишлисте
	wishlistRouter.HandleFunc("/{wishlist_id}/items", wishlistHandlers.AddItem).Methods("POST")
	wishlistRouter.HandleFunc("/{wishlist_id}/items", wishlistHandlers.ListItems).Methods("GET")

	// айтемы
	itemHandlers := item.NewHandlers(log, wishlistServiceClient)
	itemsRouter := apiRouter.PathPrefix("/items").Subrouter()
	itemsRouter.HandleFunc("/{id}", itemHandlers.GetItem).Methods("GET")
	itemsRouter.HandleFunc("/{id}", itemHandlers.UpdateItem).Methods("PUT")
	itemsRouter.HandleFunc("/{id}", itemHandlers.DeleteItem).Methods("DELETE")
	// бронирование в вишлисте
	itemsRouter.HandleFunc("/{id}/book", itemHandlers.BookItem).Methods("POST")
	itemsRouter.HandleFunc("/{id}/unbook", itemHandlers.UnbookItem).Methods("POST")

	//бронирование
	bookingHandlers := booking.NewHandlers(log, wishlistServiceClient)
	apiRouter.HandleFunc("/bookings", bookingHandlers.GetUserBookings).Methods("GET")

	// swagger ui
	if cfg.SwaggerConfig.NeedRuning {
		router.PathPrefix(cfg.SwaggerConfig.URL).Handler(httpSwagger.WrapHandler)
	}

	server := &http.Server{
		Addr:         cfg.ServerConfig.Addres + ":" + cfg.ServerConfig.Port,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.Timeout,
	}
	application.server = server

	return &application
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

	a.authServiceClient.Close()
	a.wishlistServiceClient.Close()

	log.Info("Stop server end")

}
