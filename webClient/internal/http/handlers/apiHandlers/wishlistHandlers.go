package apihandlers

import (
	"encoding/json"
	"httpClient/internal/http/clients"
	"httpClient/internal/http/handlers"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type WishlistHandlers struct {
	log                   *slog.Logger
	wishlistServiceClient clients.IWishlistGatewayService
}

func NewWishlistHandlers(log *slog.Logger, client clients.IWishlistGatewayService) *WishlistHandlers {
	return &WishlistHandlers{
		log:                   log,
		wishlistServiceClient: client,
	}
}

func (h *WishlistHandlers) CreateWishlist(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.wishlist.CreateWishlist"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	token, err := h.getTokenFromRequest(r)
	if err != nil {
		handlers.SendError(log, w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	var req handlers.CreateWishlistRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendError(log, w, http.StatusBadRequest, "invalid request body", "invalid request body: "+err.Error())
		return
	}

	if req.Title == "" {
		handlers.SendError(log, w, http.StatusBadRequest, "title is required", "title is required")
		return
	}

	resp, err := h.wishlistServiceClient.CreateWishlist(token, req)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "create wishlist failed", "create wishlist failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

func (h *WishlistHandlers) GetPublicWishlists(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.wishlist.GetPublicWishlists"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	token, err := h.getTokenFromRequest(r)
	if err != nil {
		handlers.SendError(log, w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	// Параметры пагинации
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	resp, err := h.wishlistServiceClient.GetPublicWishlists(token, int32(page), int32(limit))
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "get wishlist failed", "get wishlist failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}

}

func (h *WishlistHandlers) GetUserWishlists(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.wishlist.GetUserWishlists"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	token, err := h.getTokenFromRequest(r)
	if err != nil {
		handlers.SendError(log, w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Параметры пагинации
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	resp, err := h.wishlistServiceClient.GetUserWishlists(token, userID, int32(page), int32(limit))
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "get wishlist failed", "get wishlist failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

func (h *WishlistHandlers) GetWishlist(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.wishlist.GetWishlist"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	token, err := h.getTokenFromRequest(r)
	if err != nil {
		handlers.SendError(log, w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	vars := mux.Vars(r)
	wishlistID := vars["id"]

	resp, err := h.wishlistServiceClient.GetWishlist(token, wishlistID)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "get wishlist failed", "get wishlist failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

func (h *WishlistHandlers) UpdateWishlist(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.wishlist.UpdateWishlist"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	token, err := h.getTokenFromRequest(r)
	if err != nil {
		handlers.SendError(log, w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	vars := mux.Vars(r)
	wishlistID := vars["id"]

	var req handlers.UpdateWishlistRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendError(log, w, http.StatusBadRequest, "invalid request body", "invalid request body: "+err.Error())
		return
	}

	if req.Title == "" {
		handlers.SendError(log, w, http.StatusBadRequest, "title is required", "title is required")
		return
	}

	resp, err := h.wishlistServiceClient.UpdateWishlist(token, wishlistID, req)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "update wishlist failed", "update wishlist failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}

}
func (h *WishlistHandlers) DeleteWishlist(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.wishlist.DeleteWishlist"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	token, err := h.getTokenFromRequest(r)
	if err != nil {
		handlers.SendError(log, w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	vars := mux.Vars(r)
	wishlistID := vars["id"]

	err = h.wishlistServiceClient.DeleteWishlist(token, wishlistID)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "update wishlist failed", "update wishlist failed: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WishlistHandlers) getTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", http.ErrNoCookie
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", http.ErrNoCookie
	}

	return parts[1], nil
}
