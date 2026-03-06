package middleware

import (
	"net/http"

	"github.com/asgwg01/wishlists/pkg/types/trace"
)

func TraceMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Пытаемся получить trace ID из заголовка
			traceID := r.Header.Get(trace.HeaderKey)

			// Если нет - генерируем новый
			if traceID == "" {
				traceID = trace.GenerateTraceID()
			}

			// Добавляем trace ID в ответ (чтобы клиент мог видеть)
			w.Header().Set(trace.HeaderKey, traceID)

			// Добавляем в контекст
			ctx := trace.WithTraceID(r.Context(), traceID)

			// Передаем дальше
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
