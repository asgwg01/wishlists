package middleware

import (
	"httpClient/internal/config"
	"httpClient/internal/http/handlers"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

//var store = sessions.NewCookieStore([]byte("your-secret-key"))

// func init(cfg *config.Config) {
// 	store = sessions.NewCookieStore([]byte("your-secret-key"))
// 	store.Options = &sessions.Options{
// 		Path:     "/",
// 		MaxAge:   86400 * 7,
// 		HttpOnly: true,
// 		SameSite: http.SameSiteLaxMode,
// 		Secure:   false, // true в production с HTTPS
// 	}
// }

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserTokenKey contextKey = "user_token"
)

func AuthMiddleware(log *slog.Logger, cfg *config.ServerConfig) func(http.Handler) http.Handler {

	store := sessions.NewCookieStore([]byte(cfg.Secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, //false,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const logPrefix = "http.middleware.Auth"
			log := log.With(
				slog.String("where", logPrefix),
			)

			log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

			// Проверяем наличие токена
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handlers.SendError(log, w, http.StatusUnauthorized, "unauthorized", "unauthorized")
				return
			}

			// Проверяем формат токена
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				handlers.SendError(log, w, http.StatusUnauthorized, "Invalid authorization header", "Invalid authorization header")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
