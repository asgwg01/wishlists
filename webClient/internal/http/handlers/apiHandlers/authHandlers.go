package apihandlers

import (
	"encoding/json"
	"httpClient/internal/http/clients"
	"httpClient/internal/http/handlers"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type AuthHandlers struct {
	log                   *slog.Logger
	wishlistServiceClient clients.IWishlistGatewayService
}

func NewAuthHandlers(log *slog.Logger, client clients.IWishlistGatewayService) *AuthHandlers {
	return &AuthHandlers{
		log:                   log,
		wishlistServiceClient: client,
	}
}

func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.auth.Register"
	log := h.log.With(
		slog.String("where", logPrefix),
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

	resp, err := h.wishlistServiceClient.Register(req)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.auth.Login"
	log := h.log.With(
		slog.String("where", logPrefix),
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

	resp, err := h.wishlistServiceClient.Login(req)
	if err != nil {
		handlers.SendError(log, w, http.StatusUnauthorized, "login failed", "login failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

func (h *AuthHandlers) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.auth.GetCurrentUser"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		handlers.SendError(log, w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		handlers.SendError(log, w, http.StatusUnauthorized, "Invalid authorization header", "Invalid authorization header")
		return
	}
	token := parts[1]

	resp, err := h.wishlistServiceClient.GetCurrentUser(token)
	if err != nil {
		handlers.SendError(log, w, http.StatusUnauthorized, "login failed", "login failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

func (h *AuthHandlers) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.auth.GetUserInfo"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		handlers.SendError(log, w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		handlers.SendError(log, w, http.StatusUnauthorized, "Invalid authorization header", "Invalid authorization header")
		return
	}
	token := parts[1]

	vars := mux.Vars(r)
	userID := vars["id"]

	resp, err := h.wishlistServiceClient.GetUserInfo(token, userID)
	if err != nil {
		handlers.SendError(log, w, http.StatusUnauthorized, "get user info failed", "get user info  failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}
