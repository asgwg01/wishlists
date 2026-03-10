package apihandlers

import (
	"encoding/json"
	"httpClient/internal/http/clients"
	"httpClient/internal/http/handlers"
	"log/slog"
	"net/http"
	"strings"
)

type BookingHandlers struct {
	log                   *slog.Logger
	wishlistServiceClient clients.IWishlistGatewayService
}

func NewBookingHandlers(log *slog.Logger, client clients.IWishlistGatewayService) *BookingHandlers {
	return &BookingHandlers{
		log:                   log,
		wishlistServiceClient: client,
	}
}

func (h *BookingHandlers) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.booking.GetUserBookings"
	log := h.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	token, err := h.getTokenFromRequest(r)
	if err != nil {
		handlers.SendError(log, w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	resp, err := h.wishlistServiceClient.GetUserBookings(token)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "get usser bookings failed", "get usser bookings failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

func (h *BookingHandlers) getTokenFromRequest(r *http.Request) (string, error) {
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
