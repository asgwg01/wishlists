package item

import (
	"encoding/json"
	"gateway/internal/grpc/client"
	"gateway/internal/http/handlers"
	"gateway/internal/http/middleware"
	"log/slog"
	"net/http"
	"pkg/types/trace"

	"github.com/gorilla/mux"
)

type ItemHandlers struct {
	log            *slog.Logger
	wishlistClient *client.WishlistClient
}

func NewHandlers(log *slog.Logger, wishlistClient *client.WishlistClient) *ItemHandlers {
	return &ItemHandlers{
		log:            log,
		wishlistClient: wishlistClient,
	}
}

// GetItem возвращает предмет по ID
// @Summary      Получить предмет
// @Description  Возвращает предмет по его ID
// @Tags         Items
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID предмета"
// @Success      200  {object}  handlers.ItemDTO  "Предмет найден"
// @Failure      401  {object}  handlers.ErrorDTO "Не авторизован"
// @Failure      500  {object}  handlers.ErrorDTO "Внутренняя ошибка сервера"
// @Router       /items/{id} [get]
func (h *ItemHandlers) GetItem(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.item.GetItem"
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

	vars := mux.Vars(r)
	itemID := vars["id"]

	item, err := h.wishlistClient.GetItem(r.Context(), userID, itemID)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "failed to get item", "failed to get item: "+err.Error())
		return
	}

	response := handlers.ItemDTO{
		ID:          item.Id,
		WishlistID:  item.WishlistId,
		Name:        item.Name,
		Description: item.Description,
		ImageURL:    item.ImageUrl,
		ProductURL:  item.ProductUrl,
		Price:       item.Price,
		BookedBy:    item.BookedBy,
		BookedAt:    item.BookedAt,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

// UpdateItem обновляет предмет
// @Summary      Обновить предмет
// @Description  Обновляет существующий предмет (только владелец вишлиста)
// @Tags         Items
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path string                   true "ID предмета"
// @Param        request body handlers.UpdateItemRequestDTO    true "Новые данные"
// @Success      200  {object}  handlers.ItemDTO  "Предмет обновлен"
// @Failure      400  {object}  handlers.ErrorDTO "Неверный запрос"
// @Failure      401  {object}  handlers.ErrorDTO "Не авторизован"
// @Failure      500  {object}  handlers.ErrorDTO "Внутренняя ошибка сервера"
// @Router       /items/{id} [put]
func (h *ItemHandlers) UpdateItem(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.item.UpdateItem"
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

	vars := mux.Vars(r)
	itemID := vars["id"]

	var req handlers.UpdateItemRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendError(log, w, http.StatusBadRequest, "invalid request body", "invalid request body: "+err.Error())
		return
	}

	if req.Name == "" {
		handlers.SendError(log, w, http.StatusBadRequest, "name is required", "name is required")
		return
	}

	item, err := h.wishlistClient.UpdateItem(
		r.Context(),
		userID,
		userEmail,
		userName,
		itemID,
		req.Name,
		req.Description,
		req.ImageURL,
		req.ProductURL,
		req.Price,
	)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "failed to update item", "failed to update item: "+err.Error())
		return
	}

	response := handlers.ItemDTO{
		ID:          item.Id,
		WishlistID:  item.WishlistId,
		Name:        item.Name,
		Description: item.Description,
		ImageURL:    item.ImageUrl,
		ProductURL:  item.ProductUrl,
		Price:       item.Price,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

// DeleteItem удаляет предмет
// @Summary      Удалить предмет
// @Description  Удаляет предмет из вишлиста (только владелец)
// @Tags         Items
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID предмета"
// @Success      204  "Успешное удаление"
// @Failure      401  {object}  handlers.ErrorDTO "Не авторизован"
// @Failure      500  {object}  handlers.ErrorDTO "Внутренняя ошибка сервера"
// @Router       /items/{id} [delete]
func (h *ItemHandlers) DeleteItem(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.item.UpdateItem"
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

	vars := mux.Vars(r)
	itemID := vars["id"]

	err := h.wishlistClient.DeleteItem(r.Context(), userID, itemID)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "failed to delete item", "failed to delete item: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// BookItem бронирует предмет
// @Summary      Забронировать предмет
// @Description  Пользователь бронирует предмет в чужом вишлисте
// @Tags         Items
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID предмета"
// @Success      200  {object}  handlers.ItemDTO  "Предмет забронирован"
// @Failure      401  {object}  handlers.ErrorDTO "Не авторизован"
// @Failure      500  {object}  handlers.ErrorDTO "Внутренняя ошибка сервера"
// @Router       /items/{id}/book [post]
func (h *ItemHandlers) BookItem(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.booking.BookItem"
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

	vars := mux.Vars(r)
	itemID := vars["id"]

	item, err := h.wishlistClient.BookItem(r.Context(), userID, userEmail, userName, itemID)
	if err != nil {
		handlers.SendError(log, w, http.StatusInternalServerError, "failed to book item", "failed to book item: "+err.Error())
		return
	}

	response := handlers.ItemDTO{
		ID:          item.Id,
		WishlistID:  item.WishlistId,
		Name:        item.Name,
		Description: item.Description,
		ImageURL:    item.ImageUrl,
		ProductURL:  item.ProductUrl,
		Price:       item.Price,
		BookedBy:    item.BookedBy,
		BookedAt:    item.BookedAt,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

// UnbookItem отменяет бронирование
// @Summary      Отменить бронирование
// @Description  Отменяет бронирование предмета (бронировавший или владелец)
// @Tags         Items
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID предмета"
// @Success      200  {object}  handlers.ItemDTO  "Бронирование отменено"
// @Failure      401  {object}  handlers.ErrorDTO "Не авторизован"
// @Failure      500  {object}  handlers.ErrorDTO "Внутренняя ошибка сервера""
// @Router       /items/{id}/unbook [post]
func (h *ItemHandlers) UnbookItem(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.booking.UnbookItem"
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

	vars := mux.Vars(r)
	itemID := vars["id"]

	item, err := h.wishlistClient.UnbookItem(r.Context(), userID, userEmail, userName, itemID)
	if err != nil {

		handlers.SendError(log, w, http.StatusInternalServerError, "failed to unbook item", "failed to unbook item: "+err.Error())
		return
	}

	response := handlers.ItemDTO{
		ID:          item.Id,
		WishlistID:  item.WishlistId,
		Name:        item.Name,
		Description: item.Description,
		ImageURL:    item.ImageUrl,
		ProductURL:  item.ProductUrl,
		Price:       item.Price,
		BookedBy:    item.BookedBy,
		BookedAt:    item.BookedAt,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}
