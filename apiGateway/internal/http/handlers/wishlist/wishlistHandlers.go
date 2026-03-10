package wishlist

import (
	"encoding/json"
	"gateway/internal/grpc/client"
	"gateway/internal/http/handlers"
	"gateway/internal/http/middleware"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/asgwg01/wishlists/pkg/types/trace"

	"github.com/gorilla/mux"
)

type WishlistHandlers struct {
	log            *slog.Logger
	wishlistClient *client.WishlistClient
}

func NewHandlers(log *slog.Logger, wishlistClient *client.WishlistClient) *WishlistHandlers {
	return &WishlistHandlers{
		log:            log,
		wishlistClient: wishlistClient,
	}
}

// CreateWishlist создает новый вишлист
// @Summary      Создать вишлист
// @Description  Создает новый вишлист для текущего пользователя
// @Tags         Wishlists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body handlers.CreateWishlistRequestDTO true "Данные вишлиста"
// @Success      201  {object}  handlers.WishlistDTO  "Вишлист создан"
// @Failure      400  {object}  handlers.ErrorDTO     "Неверный запрос"
// @Failure      401  {object}  handlers.ErrorDTO     "Не авторизован"
// @Failure      500  {object}  handlers.ErrorDTO     "Внутренняя ошибка сервера"
// @Router       /wishlists [post]
func (h *WishlistHandlers) CreateWishlist(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.wishlist.CreateWishlist"
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

	var req handlers.CreateWishlistRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendError(log, w, http.StatusBadRequest, "invalid request body", "invalid request body: "+err.Error())
		return
	}

	if req.Title == "" {
		handlers.SendError(log, w, http.StatusBadRequest, "title is required", "title is required")
		return
	}

	wishlist, err := h.wishlistClient.CreateWishlist(
		ctx,
		userID,
		userEmail,
		userName,
		req.Title,
		req.Description,
		req.IsPublic,
	)
	if err != nil {
		handlers.SendError(log, w,
			http.StatusInternalServerError,
			"failed to create wishlist",
			"failed to create wishlist: "+err.Error(),
		)
		return
	}

	response := handlers.WishlistDTO{
		ID:          wishlist.Id,
		OwnerID:     wishlist.OwnerId,
		Title:       wishlist.Title,
		Description: wishlist.Description,
		IsPublic:    wishlist.IsPublic,
		CreatedAt:   wishlist.CreatedAt,
		UpdatedAt:   wishlist.UpdatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

// GetWishlist возвращает вишлист по ID
// @Summary      Получить вишлист
// @Description  Возвращает вишлист по его ID (публичные доступны всем, приватные только владельцу)
// @Tags         Wishlists
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID вишлиста"
// @Success      200  {object}  handlers.WishlistDTO  "Вишлист найден"
// @Failure      401  {object}  handlers.ErrorDTO     "Не авторизован для приватного вишлиста"
// @Failure      404  {object}  handlers.ErrorDTO     "Вишлист не найден"
// @Failure      500  {object}  handlers.ErrorDTO     "Внутренняя ошибка сервера"
// @Router       /wishlists/{id} [get]
func (h *WishlistHandlers) GetWishlist(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.wishlist.GetWishlist"
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
	wishlistID := vars["id"]

	wishlist, err := h.wishlistClient.GetWishlist(r.Context(), userID, wishlistID)
	if err != nil {
		handlers.SendError(log, w,
			http.StatusNotFound,
			"failed to get wishlist",
			"failed to get wishlist: "+err.Error(),
		)
		return
	}

	response := handlers.WishlistDTO{
		ID:          wishlist.Id,
		OwnerID:     wishlist.OwnerId,
		Title:       wishlist.Title,
		Description: wishlist.Description,
		IsPublic:    wishlist.IsPublic,
		CreatedAt:   wishlist.CreatedAt,
		UpdatedAt:   wishlist.UpdatedAt,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

// GetUserWishlists возвращает все вишлисты пользователя
// @Summary      Получить вишлисты пользователя
// @Description  Возвращает все вишлисты указанного пользователя (публичные или все для владельца)
// @Tags         Wishlists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id         path  string  true  "ID пользователя"
// @Success      200  {array}   handlers.WishlistListDTO  "Список вишлистов"
// @Failure      401  {object}  handlers.ErrorDTO     "Не авторизован"
// @Failure      500  {object}  handlers.ErrorDTO     "Внутренняя ошибка сервера"
// @Router       /wishlists/user/{user_id} [get]
func (h *WishlistHandlers) GetUserWishlists(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.wishlist.GetUserWishlists"
	log := h.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(r.Context())),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	ctx := r.Context()

	currentUserID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		handlers.SendError(log, w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	vars := mux.Vars(r)
	targetUserID := vars["user_id"]

	wishlists, err := h.wishlistClient.GetUserWishlists(
		r.Context(),
		targetUserID,
		currentUserID,
	)
	if err != nil {
		handlers.SendError(log, w,
			http.StatusInternalServerError,
			"failed to get user wishlists",
			"failed to get user wishlists: "+err.Error(),
		)
		return
	}

	response := handlers.WishlistListDTO{
		Wishlists:  make([]handlers.WishlistDTO, len(wishlists.Wishlists)),
		TotalCount: wishlists.TotalCount,
	}

	for i, w := range wishlists.Wishlists {
		response.Wishlists[i] = handlers.WishlistDTO{
			ID:          w.Id,
			OwnerID:     w.OwnerId,
			Title:       w.Title,
			Description: w.Description,
			IsPublic:    w.IsPublic,
			CreatedAt:   w.CreatedAt,
			UpdatedAt:   w.UpdatedAt,
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

// UpdateWishlist обновляет вишлист
// @Summary      Обновить вишлист
// @Description  Обновляет существующий вишлист (только владелец)
// @Tags         Wishlists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path string                     true "ID вишлиста"
// @Param        request body handlers.UpdateWishlistRequestDTO  true "Новые данные"
// @Success      200  {object}  handlers.WishlistDTO  "Вишлист обновлен"
// @Failure      400  {object}  handlers.ErrorDTO     "Неверный запрос"
// @Failure      401  {object}  handlers.ErrorDTO     "Не авторизован"
// @Failure      500  {object}  handlers.ErrorDTO     "Внутренняя ошибка сервера"
// @Router       /wishlists/{id} [put]
func (h *WishlistHandlers) UpdateWishlist(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.wishlist.UpdateWishlist"
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
	wishlistID := vars["id"]

	var req handlers.UpdateWishlistRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendError(log, w,
			http.StatusBadRequest,
			"finvalid request body",
			"invalid request body: "+err.Error(),
		)
		return
	}

	if req.Title == "" {
		handlers.SendError(log, w,
			http.StatusBadRequest,
			"title is required",
			"title is required",
		)
		return
	}

	wishlist, err := h.wishlistClient.UpdateWishlist(
		r.Context(),
		userID,
		wishlistID,
		req.Title,
		req.Description,
		req.IsPublic,
	)
	if err != nil {
		handlers.SendError(log, w,
			http.StatusInternalServerError,
			"failed to update wishlist",
			"failed to update wishlist: "+err.Error(),
		)
		return
	}

	response := handlers.WishlistDTO{
		ID:          wishlist.Id,
		OwnerID:     wishlist.OwnerId,
		Title:       wishlist.Title,
		Description: wishlist.Description,
		IsPublic:    wishlist.IsPublic,
		CreatedAt:   wishlist.CreatedAt,
		UpdatedAt:   wishlist.UpdatedAt,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

// DeleteWishlist удаляет вишлист
// @Summary      Удалить вишлист
// @Description  Удаляет вишлист и все его предметы (только владелец)
// @Tags         Wishlists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID вишлиста"
// @Success      204  "Успешное удаление"
// @Failure      401  {object}  handlers.ErrorDTO  "Не авторизован"
// @Failure      500  {object}  handlers.ErrorDTO  "Внутренняя ошибка сервера"
// @Router       /wishlists/{id} [delete]
func (h *WishlistHandlers) DeleteWishlist(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.wishlist.DeleteWishlist"
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
	wishlistID := vars["id"]

	err := h.wishlistClient.DeleteWishlist(r.Context(), userID, wishlistID)
	if err != nil {
		handlers.SendError(log, w,
			http.StatusInternalServerError,
			"failed to delete wishlist",
			"failed to delete wishlist: "+err.Error(),
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListPublicWishlists возвращает публичные вишлисты с пагинацией
// @Summary      Список публичных вишлистов
// @Description  Возвращает список публичных вишлистов с пагинацией
// @Tags         Wishlists
// @Accept       json
// @Produce      json
// @Param        page      query int false "Номер страницы" default(1)
// @Param        page_size query int false "Размер страницы" default(20)
// @Success      200  {object}  handlers.WishlistListDTO  "Список вишлистов"
// @Failure      500  {object}  handlers.ErrorDTO         "Внутренняя ошибка сервера"
// @Router       /wishlists/public [get]
func (h *WishlistHandlers) ListPublicWishlists(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.wishlist.ListPublicWishlists"
	log := h.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(r.Context())),
	)

	log.Info("Recive message", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		handlers.SendError(log, w,
			http.StatusInternalServerError,
			"failed parse int",
			"failed parse int: "+err.Error(),
		)
		return
	}

	pageSizeStr := r.URL.Query().Get("page_size")
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		handlers.SendError(log, w,
			http.StatusInternalServerError,
			"failed parse int",
			"failed parse int: "+err.Error(),
		)
		return
	}

	wishlists, err := h.wishlistClient.ListPublicWishlists(r.Context(), int32(page), int32(pageSize))
	if err != nil {
		handlers.SendError(log, w,
			http.StatusInternalServerError,
			"failed to list public wishlists",
			"failed to list public wishlists: "+err.Error(),
		)
		return
	}

	response := handlers.WishlistListDTO{
		Wishlists:   make([]handlers.WishlistDTO, len(wishlists.Wishlists)),
		TotalCount:  wishlists.TotalCount,
		TotalPages:  wishlists.TotalPages,
		CurrentPage: int32(page),
	}

	for i, w := range wishlists.Wishlists {
		response.Wishlists[i] = handlers.WishlistDTO{
			ID:          w.Id,
			OwnerID:     w.OwnerId,
			Title:       w.Title,
			Description: w.Description,
			IsPublic:    w.IsPublic,
			CreatedAt:   w.CreatedAt,
			UpdatedAt:   w.UpdatedAt,
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

// AddItem добавляет предмет в вишлист
// @Summary      Добавить предмет
// @Description  Добавляет новый предмет в указанный вишлист (только владелец)
// @Tags         Wishlists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        wishlist_id path string                true "ID вишлиста"
// @Param        request     body handlers.CreateItemRequestDTO true "Данные предмета"
// @Success      201  {object}  handlers.ItemDTO  "Предмет создан"
// @Failure      400  {object}  handlers.ErrorDTO "Неверный запрос"
// @Failure      401  {object}  handlers.ErrorDTO "Не авторизован"
// @Failure      500  {object}  handlers.ErrorDTO "Внутренняя ошибка сервера"
// @Router       /wishlists/{wishlist_id}/items [post]
func (h *WishlistHandlers) AddItem(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.wishlist.AddItem"
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
	wishlistID := vars["wishlist_id"]

	var req handlers.CreateItemRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendError(log, w,
			http.StatusBadRequest,
			"invalid request body",
			"invalid request body: "+err.Error(),
		)
		return
	}

	if req.Name == "" {
		handlers.SendError(log, w,
			http.StatusBadRequest,
			"name is required",
			"name is required",
		)
		return
	}

	item, err := h.wishlistClient.AddItem(
		r.Context(),
		userID,
		userEmail,
		userName,
		wishlistID,
		req.Name,
		req.Description,
		req.ImageURL,
		req.ProductURL,
		req.Price,
	)
	if err != nil {
		handlers.SendError(log, w,
			http.StatusInternalServerError,
			"failed to add item",
			"failed to add item: "+err.Error(),
		)
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

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}

// ListItems возвращает предметы вишлиста
// @Summary      Список предметов
// @Description  Возвращает список предметов в указанном вишлисте с пагинацией
// @Tags         Wishlists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        wishlist_id path string true "ID вишлиста"
// @Param        page        query int    false "Номер страницы" default(1)
// @Param        page_size   query int    false "Размер страницы" default(20)
// @Success      200  {object}  handlers.ItemListDTO  "Список предметов"
// @Failure      401  {object}  handlers.ErrorDTO     "Не авторизован"
// @Failure      500  {object}  handlers.ErrorDTO     "Внутренняя ошибка сервера"
// @Router       /wishlists/{wishlist_id}/items [get]
func (h *WishlistHandlers) ListItems(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "handlers.wishlist.ListItems"
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
	wishlistID := vars["wishlist_id"]

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		handlers.SendError(log, w,
			http.StatusInternalServerError,
			"failed parse int",
			"failed parse int: "+err.Error(),
		)
		return
	}

	pageSizeStr := r.URL.Query().Get("page_size")
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		handlers.SendError(log, w,
			http.StatusInternalServerError,
			"failed parse int",
			"failed parse int: "+err.Error(),
		)
		return
	}

	items, err := h.wishlistClient.ListItems(r.Context(), userID, wishlistID, int32(page), int32(pageSize))
	if err != nil {
		handlers.SendError(log, w,
			http.StatusInternalServerError,
			"failed to list items",
			"failed to list items: "+err.Error(),
		)
		return
	}

	response := handlers.ItemListDTO{
		Items:       make([]handlers.ItemDTO, len(items.Items)),
		TotalCount:  items.TotalCount,
		TotalPages:  items.TotalPages,
		CurrentPage: int32(page),
	}

	for i, item := range items.Items {
		response.Items[i] = handlers.ItemDTO{
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
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}
