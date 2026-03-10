package handlers

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

type PageHandler struct {
	log   *slog.Logger
	pages map[string]string
	fs    http.Handler
}

func NewPageHandler(log *slog.Logger) *PageHandler {
	return &PageHandler{
		log: log,
		pages: map[string]string{
			"index":              "index.html",
			"login":              "login.html",
			"register":           "register.html",
			"my_wishlists":       "my_wishlists.html",
			"my_wishlist_detail": "my_wishlist_detail.html",
			"browse_wishlists":   "browse_wishlists.html",
			"wishlist_view":      "wishlist_view.html",
			"404":                "404.html",
		},
		fs: http.FileServer(http.Dir("./web")),
	}
}

// serveFile - вспомогательная функция для безопасной отдачи файлов
func (h *PageHandler) serveFile(w http.ResponseWriter, r *http.Request, pageName string) {
	const logPrefix = "handlers.pages.serveFile"
	log := h.log.With(
		slog.String("where", logPrefix),
		slog.String("page_name", pageName),
	)

	log.Info("serve file", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	// Получаем имя файла из мапы
	filename, ok := h.pages[pageName]
	if !ok {
		log.Error("Page not found")
		h.NotFound(w, r)
		return
	}

	// Формируем полный путь к файлу
	fullPath := filepath.Join("./web", filename)

	// Проверяем существование файла
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Error("File not found",
			slog.String("full_path", fullPath),
			slog.String("err", err.Error()),
		)
		h.NotFound(w, r)
		return
	}

	// Устанавливаем правильный Content-Type
	if strings.HasSuffix(filename, ".html") {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}

	// Отдаем файл
	http.ServeFile(w, r, fullPath)
}

// Публичные страницы
func (h *PageHandler) Index(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.pages.Index"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	// Проверяем, что запрос именно к корню
	if r.URL.Path != "/" {
		h.NotFound(w, r)
		return
	}
	h.serveFile(w, r, "index")
}

func (h *PageHandler) Login(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.pages.Login"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	h.serveFile(w, r, "login")
}

func (h *PageHandler) Register(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.pages.Register"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	h.serveFile(w, r, "register")
}

// Защищенные страницы (требуют авторизации)
func (h *PageHandler) MyWishlists(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.pages.MyWishlists"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	h.serveFile(w, r, "my_wishlists")
}

func (h *PageHandler) MyWishlistDetail(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.pages.MyWishlistDetail"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Redirect(w, r, "/my_wishlists", http.StatusSeeOther)
		return
	}

	h.serveFile(w, r, "my_wishlist_detail")
}

func (h *PageHandler) BrowseWishlists(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.pages.BrowseWishlists"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	h.serveFile(w, r, "browse_wishlists")
}

func (h *PageHandler) WishlistView(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.pages.WishlistView"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Redirect(w, r, "/browse_wishlists", http.StatusSeeOther)
		return
	}

	h.serveFile(w, r, "wishlist_view")
}

// Обработчик 404
func (h *PageHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.pages.NotFound"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	// Пробуем отдать кастомную 404 страницу
	notFoundPath := filepath.Join("./web", "404.html")
	if _, err := os.Stat(notFoundPath); err == nil {
		http.ServeFile(w, r, notFoundPath)
		return
	}

	// Если нет кастомной 404, отдаем простой текст
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <title>404 - Страница не найдена</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #0f0f0f;
            color: #f3f4f6;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            text-align: center;
        }
        .error-container {
            max-width: 600px;
            padding: 2rem;
        }
        h1 {
            font-size: 4rem;
            margin-bottom: 1rem;
            background: linear-gradient(135deg, #7c3aed, #10b981);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }
        p {
            font-size: 1.2rem;
            color: #9ca3af;
            margin-bottom: 2rem;
        }
        a {
            display: inline-block;
            padding: 0.75rem 1.5rem;
            background: linear-gradient(135deg, #7c3aed, #6d28d9);
            color: white;
            text-decoration: none;
            border-radius: 0.5rem;
            transition: all 0.2s;
        }
        a:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 6px rgba(0,0,0,0.3);
        }
    </style>
</head>
<body>
    <div class="error-container">
        <h1>404</h1>
        <p>Страница, которую вы ищете, не существует или была перемещена.</p>
        <a href="/">Вернуться на главную</a>
    </div>
</body>
</html>`))
}

// StaticFiles обработчик для статических файлов
func (h *PageHandler) StaticFiles() http.Handler {
	const logPrefix = "handlers.pages.StaticFiles"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Static Files")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Предотвращаем directory traversal атаки
		requestedPath := r.URL.Path
		if strings.Contains(requestedPath, "..") {
			SendError(log, w, http.StatusForbidden, "Forbidden", "Forbidden")
			return
		}

		// Устанавливаем правильные заголовки для статических файлов
		switch {
		case strings.HasSuffix(requestedPath, ".css"):
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case strings.HasSuffix(requestedPath, ".js"):
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		case strings.HasSuffix(requestedPath, ".png"):
			w.Header().Set("Content-Type", "image/png")
		case strings.HasSuffix(requestedPath, ".jpg"), strings.HasSuffix(requestedPath, ".jpeg"):
			w.Header().Set("Content-Type", "image/jpeg")
		case strings.HasSuffix(requestedPath, ".svg"):
			w.Header().Set("Content-Type", "image/svg+xml")
		case strings.HasSuffix(requestedPath, ".webp"):
			w.Header().Set("Content-Type", "image/webp")
		case strings.HasSuffix(requestedPath, ".ico"):
			w.Header().Set("Content-Type", "image/x-icon")
		}

		// Добавляем кэширование для статических файлов
		w.Header().Set("Cache-Control", "public, max-age=86400")

		h.fs.ServeHTTP(w, r)
	})
}
