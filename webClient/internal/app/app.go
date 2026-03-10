package app

import (
	"context"
	"errors"
	"httpClient/internal/config"
	wishlistsgateway "httpClient/internal/http/clients/wishlistsGateway"
	"httpClient/internal/http/handlers"
	apihandlers "httpClient/internal/http/handlers/apiHandlers"
	"httpClient/internal/http/middleware"

	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	log    *slog.Logger
	server *http.Server
}

func New(log *slog.Logger, cfg *config.Config) *App {

	application := App{}
	application.log = log

	// клиент сервиса вишлистов
	wlgService := wishlistsgateway.New(log, cfg)

	// хендлеры
	pageHandlers := handlers.NewPageHandler(log)
	wishlistServiceHandlers := apihandlers.NewWishlistGatewayHandlers(log, wlgService)

	// routers
	router := mux.NewRouter()

	// static
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js"))))

	// html
	// public без авторизации
	router.HandleFunc("/", pageHandlers.Index).Methods("GET")
	router.HandleFunc("/login", pageHandlers.Login).Methods("GET")
	router.HandleFunc("/register", pageHandlers.Register).Methods("GET")

	// protected
	protected := router.PathPrefix("/").Subrouter()
	protected.HandleFunc("/my_wishlists", pageHandlers.MyWishlists).Methods("GET")
	protected.HandleFunc("/my_wishlist/{id}", pageHandlers.MyWishlistDetail).Methods("GET")
	protected.HandleFunc("/browse_wishlists", pageHandlers.BrowseWishlists).Methods("GET")
	protected.HandleFunc("/wishlist_view/{id}", pageHandlers.WishlistView).Methods("GET")

	// api вызовы
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// public без авторизации
	apiRouter.HandleFunc("/auth/register", wishlistServiceHandlers.AuthHandlers.Register).Methods("POST")
	apiRouter.HandleFunc("/auth/login", wishlistServiceHandlers.AuthHandlers.Login).Methods("POST")

	// protected
	protectedAPI := apiRouter.PathPrefix("/").Subrouter()
	protectedAPI.Use(middleware.AuthMiddleware(log, &cfg.ServerConfig))
	protectedAPI.HandleFunc("/auth/self", wishlistServiceHandlers.AuthHandlers.GetCurrentUser).Methods("GET")
	protectedAPI.HandleFunc("/auth/user/{id}", wishlistServiceHandlers.AuthHandlers.GetUserInfo).Methods("GET")

	// Wishlist client ручки
	// Wishlists
	protectedAPI.HandleFunc("/wishlists", wishlistServiceHandlers.WishlistHandlers.CreateWishlist).Methods("POST")
	protectedAPI.HandleFunc("/wishlists/public", wishlistServiceHandlers.WishlistHandlers.GetPublicWishlists).Methods("GET")
	protectedAPI.HandleFunc("/wishlists/user/{user_id}", wishlistServiceHandlers.WishlistHandlers.GetUserWishlists).Methods("GET")
	protectedAPI.HandleFunc("/wishlists/{id}", wishlistServiceHandlers.WishlistHandlers.GetWishlist).Methods("GET")
	protectedAPI.HandleFunc("/wishlists/{id}", wishlistServiceHandlers.WishlistHandlers.UpdateWishlist).Methods("PUT")
	protectedAPI.HandleFunc("/wishlists/{id}", wishlistServiceHandlers.WishlistHandlers.DeleteWishlist).Methods("DELETE")

	// Items
	protectedAPI.HandleFunc("/wishlists/{wishlist_id}/items", wishlistServiceHandlers.ItemHandlers.CreateItem).Methods("POST")
	protectedAPI.HandleFunc("/wishlists/{wishlist_id}/items", wishlistServiceHandlers.ItemHandlers.GetItems).Methods("GET")
	protectedAPI.HandleFunc("/items/{id}", wishlistServiceHandlers.ItemHandlers.GetItem).Methods("GET")
	protectedAPI.HandleFunc("/items/{id}", wishlistServiceHandlers.ItemHandlers.UpdateItem).Methods("PUT")
	protectedAPI.HandleFunc("/items/{id}", wishlistServiceHandlers.ItemHandlers.DeleteItem).Methods("DELETE")
	protectedAPI.HandleFunc("/items/{id}/book", wishlistServiceHandlers.ItemHandlers.BookItem).Methods("POST")
	protectedAPI.HandleFunc("/items/{id}/unbook", wishlistServiceHandlers.ItemHandlers.UnbookItem).Methods("POST")

	// Bookings
	protectedAPI.HandleFunc("/bookings", wishlistServiceHandlers.BookingHandlers.GetUserBookings).Methods("GET")

	// 404 общий
	router.NotFoundHandler = http.HandlerFunc(pageHandlers.NotFound)

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

	log.Info("Stop server end")

}
