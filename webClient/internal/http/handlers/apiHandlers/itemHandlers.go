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

type ItemHandlers struct {
	log                   *slog.Logger
	wishlistServiceClient clients.IWishlistGatewayService
}

func NewItemHandlers(log *slog.Logger, client clients.IWishlistGatewayService) *ItemHandlers {
	return &ItemHandlers{
		log:                   log,
		wishlistServiceClient: client,
	}
}

func (h *ItemHandlers) CreateItem(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.item.CreateItem"
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
	wishlistID := vars["wishlist_id"]

	var req handlers.CreateItemRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendError(log, w, http.StatusBadRequest, "invalid request body", "invalid request body: "+err.Error())
		return
	}

	if req.Name == "" {
		handlers.SendError(log, w, http.StatusBadRequest, "name is required", "title is required")
		return
	}

	resp, err := h.wishlistServiceClient.CreateItem(token, wishlistID, req)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "create item failed", "create item failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

func (h *ItemHandlers) GetItems(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.item.GetItems"
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
	wishlistID := vars["wishlist_id"]

	// Параметры пагинации
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	resp, err := h.wishlistServiceClient.GetItems(token, wishlistID, int32(page), int32(limit))
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "get wishlist items failed", "get wishlist items failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}
func (h *ItemHandlers) GetItem(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.item.GetItem"
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
	itemID := vars["id"]

	resp, err := h.wishlistServiceClient.GetItem(token, itemID)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "get item failed", "get item failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}
func (h *ItemHandlers) UpdateItem(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.item.UpdateItem"
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
	itemID := vars["id"]

	var req handlers.UpdateItemRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendError(log, w, http.StatusBadRequest, "invalid request body", "invalid request body: "+err.Error())
		return
	}

	if req.Name == "" {
		handlers.SendError(log, w, http.StatusBadRequest, "name is required", "title is required")
		return
	}

	resp, err := h.wishlistServiceClient.UpdateItem(token, itemID, req)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "update item failed", "update item failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}
func (h *ItemHandlers) DeleteItem(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.item.DeleteItem"
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
	itemID := vars["id"]

	err = h.wishlistServiceClient.DeleteItem(token, itemID)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "delete item failed", "delete item failed: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
func (h *ItemHandlers) BookItem(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.item.BookItem"
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
	itemID := vars["id"]

	resp, err := h.wishlistServiceClient.BookItem(token, itemID)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "book item failed", "book item failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}
func (h *ItemHandlers) UnbookItem(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.item.UnbookItem"
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
	itemID := vars["id"]

	resp, err := h.wishlistServiceClient.UnbookItem(token, itemID)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "unbook item failed", "unbook item failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

func (h *ItemHandlers) getTokenFromRequest(r *http.Request) (string, error) {
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
