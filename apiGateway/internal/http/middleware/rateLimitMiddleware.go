package middleware

import (
	"gateway/internal/config"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type RateLimitter struct {
	ipMap          map[string][]time.Time
	mtx            sync.Mutex
	limit          int
	windowDuration time.Duration
}

func NewRateLimitter(cfg config.RateLimitsConfig) *RateLimitter {
	return &RateLimitter{
		limit:          cfg.Limit,
		windowDuration: cfg.Duration,
		ipMap:          make(map[string][]time.Time),
	}
}

func (rl *RateLimitter) allow(ip string) bool {
	rl.mtx.Lock()
	defer rl.mtx.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.windowDuration)

	// Получаем историю запросов для IP
	requests, exists := rl.ipMap[ip]

	// Если первый запрос - создаем запись
	if !exists {
		rl.ipMap[ip] = []time.Time{now}
		return true
	}

	// Фильтруем только запросы внутри текущего окна
	valid := make([]time.Time, 0)
	for _, t := range requests {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}

	// Проверяем лимит
	if len(valid) >= rl.limit {
		rl.ipMap[ip] = valid
		return false // Превышен лимит
	}

	// Добавляем новый запрос
	valid = append(valid, now)
	rl.ipMap[ip] = valid
	return true
}

func RateLimitMiddleware(log *slog.Logger, limitter *RateLimitter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const logPrefix = "http.middleware.rateLimit"
			log := log.With(
				slog.String("where", logPrefix),
			)

			ip := r.RemoteAddr

			// Проверяем лимит
			if !limitter.allow(ip) {
				strconv.Itoa(limitter.limit)
				// Добавляем заголовки с информацией о лимитах
				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limitter.limit))
				w.Header().Set("X-RateLimit-Reset", strconv.Itoa(int(limitter.windowDuration.Seconds())))

				// Логируем превышение
				log.Warn("Rate limit exceeded", slog.String("ip", ip))

				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}

			// Пропускаем запрос дальше
			next.ServeHTTP(w, r)
		})
	}
}
