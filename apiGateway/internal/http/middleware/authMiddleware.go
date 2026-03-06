package middleware

import (
	"context"
	"gateway/internal/grpc/client"
	"log/slog"
	"net/http"
	"strings"
)

const (
	UserIDKey    string = "user_id"
	UserEmailKey string = "user_email"
	UserNameKey  string = "user_name"
)

// AuthMiddleware проверяет JWT токен и добавляет информацию о пользователе в контекст
func AuthMiddleware(log *slog.Logger, authClient *client.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Пропускаем публичные маршруты
			if isPublicPath(r.URL.Path) {
				log.Info("public url", slog.String("url", r.URL.Path))
				next.ServeHTTP(w, r)
				return
			}

			// Получаем токен из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Error("missing authorization header", slog.String("url", r.URL.Path))
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			// Проверяем формат "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// Валидируем токен через Auth Service
			userInfo, err := authClient.ValidateToken(r.Context(), token)
			if err != nil {
				http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Добавляем информацию о пользователе в контекст
			ctx := context.WithValue(r.Context(), UserIDKey, userInfo.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, userInfo.Email)
			ctx = context.WithValue(ctx, UserNameKey, userInfo.Name)

			// Передаем управление дальше с обновленным контекстом
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// isPublicPath проверяет, является ли путь публичным (не требует аутентификации)
func isPublicPath(path string) bool {
	publicPaths := []string{
		"/health",
		"/ready",
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/swagger",
	}

	for _, p := range publicPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
