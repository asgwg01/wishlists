package auth

import (
	"encoding/json"
	"gateway/internal/grpc/client"
	"gateway/internal/http/handlers"
	"gateway/internal/http/middleware"
	"log/slog"
	"net/http"

	"github.com/asgwg01/wishlists/pkg/types/trace"
)

type AuthHandlers struct {
	log        *slog.Logger
	authClient *client.AuthClient
}

func NewHandlers(log *slog.Logger, authClient *client.AuthClient) *AuthHandlers {
	return &AuthHandlers{
		log:        log,
		authClient: authClient,
	}
}

// Register регистрирует нового пользователя
// @Summary      Регистрация пользователя
// @Description  Создает нового пользователя в системе
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body handlers.RegisterRequestDTO true "Данные для регистрации"
// @Success      201  {object}  handlers.AuthDTO  "Успешная регистрация"
// @Failure      400  {object}  handlers.ErrorDTO "Неверный запрос"
// @Failure      409  {object}  handlers.ErrorDTO "Пользователь уже существует"
// @Failure      500  {object}  handlers.ErrorDTO "Внутренняя ошибка сервера"
// @Router       /auth/register [post]
func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.auth.Register"
	log := h.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(r.Context())),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	var req handlers.RegisterRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendError(log, w, http.StatusBadRequest, "invalid request body", "invalid request body: "+err.Error())
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		handlers.SendError(log, w, http.StatusBadRequest, "email, password and name are required", "email, password and name are required")
		return
	}

	userInfo, token, err := h.authClient.Register(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "registration failed", "registration failed: "+err.Error())
		return
	}

	response := handlers.AuthDTO{
		Token:  token,
		UserID: userInfo.UserID,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

// Login аутентифицирует пользователя
// @Summary      Вход в систему
// @Description  Аутентификация пользователя и получение JWT токена
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body handlers.LoginRequestDTO true "Учетные данные"
// @Success      200  {object}  handlers.AuthDTO  "Успешный вход"
// @Failure      400  {object}  handlers.ErrorDTO "Неверный запрос"
// @Failure      401  {object}  handlers.ErrorDTO "Неверные учетные данные"
// @Failure      500  {object}  handlers.ErrorDTO "Внутренняя ошибка сервера"
// @Router       /auth/login [post]
func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.auth.Login"
	log := h.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(r.Context())),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	var req handlers.LoginRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendError(log, w, http.StatusBadRequest, "invalid request body", "invalid request body: "+err.Error())
		return
	}

	if req.Email == "" || req.Password == "" {
		handlers.SendError(log, w, http.StatusBadRequest, "email and password are required", "email and password are required")
		return
	}

	userInfo, token, err := h.authClient.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		handlers.SendError(log, w, http.StatusUnauthorized, "login failed", "login failed: "+err.Error())
		return
	}

	response := handlers.AuthDTO{
		Token:  token,
		UserID: userInfo.UserID,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

// GetCurrentUser возвращает информацию о текущем пользователе
// @Summary      Информация о текущем пользователе
// @Description  Возвращает данные аутентифицированного пользователя
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Success      200  {object}  handlers.UserInfoDTO  "Информация о пользователе"
// @Failure      401  {object}  handlers.ErrorDTO     "Не авторизован"
// @Failure      500  {object}  handlers.ErrorDTO     "Внутренняя ошибка сервера"
// @Router       /auth/self [get]
func (h *AuthHandlers) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.auth.GetCurrentUser"
	log := h.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(r.Context())),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	ctx := r.Context()

	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		handlers.SendError(log, w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	userEmail, _ := ctx.Value(middleware.UserEmailKey).(string)
	userName, _ := ctx.Value(middleware.UserNameKey).(string)

	response := handlers.UserInfoDTO{
		UserID: userID,
		Email:  userEmail,
		Name:   userName,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}
