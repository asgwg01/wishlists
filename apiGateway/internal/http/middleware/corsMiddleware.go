package middleware

import (
	"gateway/internal/config"
	"net/http"
	"strings"
)

func CORSMiddleware(cfg config.CORSConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const logPrefix = "http.middleware.logging"
			origin := r.Header.Get("Origin")

			// Проверяем разрешен ли origin
			allowed := false
			for _, o := range cfg.Origins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}

			if allowed && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.Methods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.Headers, ", "))

				if cfg.Credendials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}

			// Обрабатываем preflight запросы
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)

		})
	}
}
