package booking

import (
	"encoding/json"
	"gateway/internal/grpc/client"
	"gateway/internal/http/handlers"
	"gateway/internal/http/middleware"
	"log/slog"
	"net/http"

	"github.com/asgwg01/wishlists/pkg/types/trace"
)

type BookingHandlers struct {
	log            *slog.Logger
	wishlistClient *client.WishlistClient
}

func NewHandlers(log *slog.Logger, wishlistClient *client.WishlistClient) *BookingHandlers {
	return &BookingHandlers{
		log:            log,
		wishlistClient: wishlistClient,
	}
}

// GetUserBookings возвращает все бронирования пользователя
// @Summary      Бронирования пользователя
// @Description  Возвращает все предметы, забронированные текущим пользователем
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  handlers.BookingListDTO  "Список бронирований"
// @Failure      401  {object}  handlers.ErrorDTO     "Не авторизован"
// @Failure      500  {object}  handlers.ErrorDTO     "Внутренняя ошибка сервера"
// @Router       /bookings [get]
func (h *BookingHandlers) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.booking.GetUserBookings"
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

	bookings, err := h.wishlistClient.GetUserBookings(r.Context(), userID)
	if err != nil {
		handlers.SendError(log, w,
			http.StatusInternalServerError,
			"failed to get user bookings",
			"failed to get user bookings: "+err.Error(),
		)
		return
	}

	response := handlers.BookingListDTO{
		Bookings:   make([]handlers.BookingDTO, len(bookings.Bookings)),
		TotalCount: bookings.TotalCount,
	}

	for i, b := range bookings.Bookings {
		response.Bookings[i] = handlers.BookingDTO{
			ItemID:     b.ItemId,
			WishlistID: b.WishlistId,
			UserID:     b.UserId,
			BookedAt:   b.BookedAt,
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}
