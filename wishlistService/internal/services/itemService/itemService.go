package itemservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"pkg/types/price"
	"pkg/types/trace"
	"wishlistService/internal/domain/models"
	"wishlistService/internal/domain/types"
	"wishlistService/internal/grpc/client"
	"wishlistService/internal/producer"
	"wishlistService/internal/storage"

	"github.com/google/uuid"
)

type IItemService interface {
	Add(ctx context.Context, wishlistID uuid.UUID, userID uuid.UUID, name, description, imageURL, productURL string, price price.Price) (models.Item, error)
	Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (models.Item, error)
	List(ctx context.Context, wishlistID uuid.UUID, userID uuid.UUID, page, pageSize int) ([]models.Item, int, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, name, description, imageURL, productURL string, price price.Price) (models.Item, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	Book(ctx context.Context, itemID uuid.UUID, userID uuid.UUID) (models.Item, error)
	Unbook(ctx context.Context, itemID uuid.UUID, userID uuid.UUID) (models.Item, error)
	GetUserBookings(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Item, error)
}

type ItemService struct {
	log             *slog.Logger
	wishlistStorage storage.IWishlistStorage
	itemStorage     storage.IItemStorage
	authClient      *client.AuthClient
	kafkaProd       *producer.KafkaProducer
}

func New(
	log *slog.Logger,
	wishlistStorage storage.IWishlistStorage,
	itemStorage storage.IItemStorage,
	authClient *client.AuthClient,
	kafkaProd *producer.KafkaProducer,
) *ItemService {
	return &ItemService{
		log:             log,
		itemStorage:     itemStorage,
		wishlistStorage: wishlistStorage,
		authClient:      authClient,
		kafkaProd:       kafkaProd,
	}
}

func (s *ItemService) Add(ctx context.Context,
	wishlistID uuid.UUID,
	userID uuid.UUID,
	name string,
	description string,
	imageURL string,
	productURL string,
	price price.Price,
) (models.Item, error) {
	const logPrefix = "service.itemService.Add"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("wishlistID", wishlistID.String()),
		slog.String("userId", userID.String()),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Add item to wishlist")

	wl, err := s.wishlistStorage.GetWishlistByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, storage.ErrorWishlistNotExist) {
			log.Error("Wishlist not found", slog.String("err", err.Error()))
			return models.Item{}, types.ErrorWishlistNotFound
		}
		log.Error("Error get wishlist", slog.String("err", err.Error()))
		return models.Item{}, err
	}

	if !wl.CanEdit(userID) {
		log.Error("Can not edit wishlist")
		return models.Item{}, types.ErrorAccessDenied
	}

	newItem := models.NewItem(wishlistID, name, description, imageURL, productURL, price)

	newItem, err = s.itemStorage.CreateItem(ctx, newItem)
	if err != nil {
		if errors.Is(err, storage.ErrorItemAlreadyExist) {
			log.Error("Item already exists", slog.String("err", err.Error()))
			return models.Item{}, types.ErrorItemAlreadyExist
		}
		log.Error("Error create item", slog.String("err", err.Error()))
		return models.Item{}, fmt.Errorf("failed to create item: %w", err)
	}

	return newItem, nil
}
func (s *ItemService) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (models.Item, error) {
	const logPrefix = "service.itemService.Get"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("itemID", id.String()),
		slog.String("userId", userID.String()),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Get item")

	item, err := s.itemStorage.GetItemByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrorItemNotExist) {
			log.Error("Wishlist not found", slog.String("err", err.Error()))
			return models.Item{}, types.ErrorWishlistNotFound
		}
		log.Error("Error get wishlist", slog.String("err", err.Error()))
		return models.Item{}, err
	}

	wl, err := s.wishlistStorage.GetWishlistByID(ctx, item.WishlistID)
	if err != nil {
		if errors.Is(err, storage.ErrorWishlistNotExist) {
			log.Error("Wishlist not found", slog.String("err", err.Error()))
			return models.Item{}, types.ErrorWishlistNotFound
		}
		log.Error("Error get wishlist", slog.String("err", err.Error()))
		return models.Item{}, err
	}

	if !wl.CanView(userID) {
		log.Error("Can not view wishlist")
		return models.Item{}, types.ErrorAccessDenied
	}

	return item, nil
}
func (s *ItemService) List(ctx context.Context,
	wishlistID uuid.UUID,
	userID uuid.UUID,
	page int,
	pageSize int,
) ([]models.Item, int, error) {
	const logPrefix = "service.itemService.List"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.Int("page", page),
		slog.Int("pageSize", pageSize),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("List items")

	wl, err := s.wishlistStorage.GetWishlistByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, storage.ErrorWishlistNotExist) {
			log.Error("Wishlist not found", slog.String("err", err.Error()))
			return []models.Item{}, 0, types.ErrorWishlistNotFound
		}
		log.Error("Error get wishlist", slog.String("err", err.Error()))
		return []models.Item{}, 0, err
	}

	if !wl.CanView(userID) {
		log.Error("Can not view wishlist")
		return []models.Item{}, 0, types.ErrorAccessDenied
	}

	// пагинация
	offset := (page - 1) * pageSize
	items, total, err := s.itemStorage.GetItemsByWishlistID(ctx, wishlistID, pageSize, offset)
	if err != nil {
		log.Error("Can not view wishlist items", slog.String("err", err.Error()))
		return []models.Item{}, 0, fmt.Errorf("failed to get list items: %w", err)
	}

	return items, total, nil
}
func (s *ItemService) Update(ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
	name string,
	description string,
	imageURL string,
	productURL string,
	price price.Price,
) (models.Item, error) {
	const logPrefix = "service.itemService.Update"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("itemID", id.String()),
		slog.String("userId", userID.String()),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Update item")

	item, err := s.itemStorage.GetItemByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrorItemNotExist) {
			log.Error("Wishlist not found", slog.String("err", err.Error()))
			return models.Item{}, types.ErrorWishlistNotFound
		}
		log.Error("Error get wishlist", slog.String("err", err.Error()))
		return models.Item{}, err
	}

	wl, err := s.wishlistStorage.GetWishlistByID(ctx, item.WishlistID)
	if err != nil {
		if errors.Is(err, storage.ErrorWishlistNotExist) {
			log.Error("Wishlist not found", slog.String("err", err.Error()))
			return models.Item{}, types.ErrorWishlistNotFound
		}
		log.Error("Error get wishlist", slog.String("err", err.Error()))
		return models.Item{}, err
	}

	if !wl.CanEdit(userID) {
		log.Error("Can not edit wishlist")
		return models.Item{}, types.ErrorAccessDenied
	}

	if name == "" {
		log.Error("name is required")
		return models.Item{}, fmt.Errorf("name is required")
	}

	item.Update(name, description, imageURL, productURL, price)

	updatedItem, err := s.itemStorage.UpdateItem(ctx, item)
	if err != nil {
		if errors.Is(err, storage.ErrorItemNotExist) {
			return models.Item{}, types.ErrorItemNotFound
		}
		return models.Item{}, fmt.Errorf("failed to update item: %w", err)
	}

	return updatedItem, nil
}
func (s *ItemService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	const logPrefix = "service.itemService.Delete"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("itemID", id.String()),
		slog.String("userId", userID.String()),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Delete item")

	item, err := s.itemStorage.GetItemByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrorItemNotExist) {
			log.Error("Wishlist not found", slog.String("err", err.Error()))
			return types.ErrorWishlistNotFound
		}
		log.Error("Error get wishlist", slog.String("err", err.Error()))
		return err
	}

	wl, err := s.wishlistStorage.GetWishlistByID(ctx, item.WishlistID)
	if err != nil {
		if errors.Is(err, storage.ErrorWishlistNotExist) {
			log.Error("Wishlist not found", slog.String("err", err.Error()))
			return types.ErrorWishlistNotFound
		}
		log.Error("Error get wishlist", slog.String("err", err.Error()))
		return err
	}

	if !wl.CanEdit(userID) {
		log.Error("Can not edit wishlist")
		return types.ErrorAccessDenied
	}

	err = s.itemStorage.DeleteItem(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrorItemNotExist) {
			return types.ErrorItemNotFound
		}
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}
func (s *ItemService) Book(ctx context.Context, itemID uuid.UUID, userID uuid.UUID) (models.Item, error) {
	const logPrefix = "service.itemService.Book"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("itemID", itemID.String()),
		slog.String("userId", userID.String()),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Book item")

	item, err := s.itemStorage.GetItemByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, storage.ErrorItemNotExist) {
			log.Error("Wishlist not found", slog.String("err", err.Error()))
			return models.Item{}, types.ErrorItemNotFound
		}
		log.Error("Error get wishlist", slog.String("err", err.Error()))
		return models.Item{}, err
	}

	wl, err := s.wishlistStorage.GetWishlistByID(ctx, item.WishlistID)
	if err != nil {
		if errors.Is(err, storage.ErrorWishlistNotExist) {
			log.Error("Wishlist not found", slog.String("err", err.Error()))
			return models.Item{}, types.ErrorWishlistNotFound
		}
		log.Error("Error get wishlist", slog.String("err", err.Error()))
		return models.Item{}, err
	}

	isSelf := wl.OwnerID == userID
	isPublic := wl.IsPublic

	if !isSelf && !isPublic {
		log.Error("Can not book private wishlist item")
		return models.Item{}, types.ErrorAccessDenied
	}

	if err := item.Book(userID); err != nil {
		if errors.Is(err, types.ErrorItemAlreadyBooked) {
			log.Error("Can not book already booked wishlist item")
			return models.Item{}, types.ErrorItemAlreadyBooked
		}
		return models.Item{}, fmt.Errorf("failed to book item: %w", err)
	}

	item, err = s.itemStorage.UpdateItemBooking(ctx, item.ID, item.BookedBy, item.BookedAt)
	if err != nil {
		if errors.Is(err, storage.ErrorItemNotExist) {
			log.Error("Wishlist not found", slog.String("err", err.Error()))
			return models.Item{}, types.ErrorItemNotFound
		}
		log.Error("Error book item", slog.String("err", err.Error()))
		return models.Item{}, err
	}

	s.SendKafkaBookedEvent(ctx, userID, &wl, &item)

	return item, nil
}
func (s *ItemService) Unbook(ctx context.Context, itemID uuid.UUID, userID uuid.UUID) (models.Item, error) {
	const logPrefix = "service.itemService.Unbook"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("itemID", itemID.String()),
		slog.String("userId", userID.String()),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Unbook item")

	item, err := s.itemStorage.GetItemByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, storage.ErrorItemNotExist) {
			log.Error("Wishlist not found", slog.String("err", err.Error()))
			return models.Item{}, types.ErrorItemNotFound
		}
		log.Error("Error get wishlist", slog.String("err", err.Error()))
		return models.Item{}, err
	}

	wl, err := s.wishlistStorage.GetWishlistByID(ctx, item.WishlistID)
	if err != nil {
		if errors.Is(err, storage.ErrorWishlistNotExist) {
			log.Error("Wishlist not found", slog.String("err", err.Error()))
			return models.Item{}, types.ErrorWishlistNotFound
		}
		log.Error("Error get wishlist", slog.String("err", err.Error()))
		return models.Item{}, err
	}

	if item.BookedBy == nil {
		log.Error("Can not unbook not booked wishlist item")
		return models.Item{}, types.ErrorItemNotBooked
	}

	isSelf := wl.OwnerID == userID
	isPublic := wl.IsPublic

	if !isSelf && !isPublic {
		log.Error("Can not unbook private wishlist item")
		return models.Item{}, types.ErrorAccessDenied
	}

	if err := item.Unbook(userID); err != nil {
		if errors.Is(err, types.ErrorItemNotBooked) {
			log.Error("Can not unbook not booked wishlist item")
			return models.Item{}, types.ErrorItemAlreadyBooked
		}
		if errors.Is(err, types.ErrorAccessDenied) {
			log.Error("Can not unbook other user booked")
			return models.Item{}, types.ErrorAccessDenied
		}
		return models.Item{}, fmt.Errorf("failed to unbook item: %w", err)
	}

	lastOwnerID := uuid.UUID{}
	if item.BookedBy != nil {
		lastOwnerID = *item.BookedBy
	}
	item, err = s.itemStorage.UpdateItemBooking(ctx, item.ID, item.BookedBy, item.BookedAt)
	if err != nil {
		if errors.Is(err, storage.ErrorItemNotExist) {
			log.Error("Wishlist not found", slog.String("err", err.Error()))
			return models.Item{}, types.ErrorItemNotFound
		}
		log.Error("Error book item", slog.String("err", err.Error()))
		return models.Item{}, err
	}

	s.SendKafkaUnbookedEvent(ctx, userID, lastOwnerID, &wl, &item)
	s.SendKafkaBookedEvent(ctx, userID, &wl, &item)

	return item, nil
}

// GetUserBookings возвращает все предметы, забронированные пользователем
func (s *ItemService) GetUserBookings(ctx context.Context, userID uuid.UUID, page int, pageSize int) ([]models.Item, error) {
	const logPrefix = "service.itemService.GetUserBookings"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("userId", userID.String()),
		slog.Int("page", page),
		slog.Int("pageSize", pageSize),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("GetUserBookings")

	// пагинация
	offset := (page - 1) * pageSize
	items, err := s.itemStorage.GetBookedItemsByUser(ctx, userID, pageSize, offset)
	if err != nil {
		log.Error("Error get booked items by user", slog.String("err", err.Error()))
		return []models.Item{}, err
	}

	return items, nil
}

func (s *ItemService) SendKafkaBookedEvent(ctx context.Context, userID uuid.UUID, wl *models.Wishlist, item *models.Item) {
	const logPrefix = "service.itemService.SendKafkaEvent"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("SendKafkaBookedEvent")

	// Получаем информацию о пользователе, который бронирует
	userInfo, err := s.authClient.GetUserInfo(ctx, userID.String())
	if err != nil {
		log.Error("Failed to get user info", slog.String("err", err.Error()))
		userInfo = &client.UserInfo{
			UserID: userID.String(),
			Email:  "unknown@email.com",
			Name:   "Unknown User",
		}
	}

	// Получаем информацию о владельце вишлиста
	ownerInfo, err := s.authClient.GetUserInfo(ctx, wl.OwnerID.String())
	if err != nil {
		log.Error("Failed to get user info", slog.String("err", err.Error()))
		ownerInfo = &client.UserInfo{
			UserID: wl.OwnerID.String(),
			Email:  "unknown@email.com",
			Name:   "Unknown Owner",
		}
	}

	evt := producer.ItemBookedEvent{
		ItemID:       item.ID.String(),
		ItemName:     item.Name,
		WishlistID:   wl.ID.String(),
		WishlistName: wl.Title,

		BookedBy:      userInfo.UserID,
		BookedByEmail: userInfo.Email,
		BookedByName:  userInfo.Name,

		OwnerID:    ownerInfo.UserID,
		OwnerEmail: ownerInfo.Email,
		OwnerName:  ownerInfo.Name,
	}

	go func() {
		if err := s.kafkaProd.PublishItemBooked(context.Background(), evt); err != nil {
			s.log.Error("Failed to publish booked event", slog.String("err", err.Error()))
		}
	}()
}

func (s *ItemService) SendKafkaUnbookedEvent(ctx context.Context, userID uuid.UUID, lastOwnerID uuid.UUID, wl *models.Wishlist, item *models.Item) {
	const logPrefix = "service.itemService.SendKafkaUnbookedEvent"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("SendKafkaUnbookedEvent")

	// Получаем информацию о пользователе, который бронирует
	userInfo, err := s.authClient.GetUserInfo(ctx, userID.String())
	if err != nil {
		log.Error("Failed to get user info", slog.String("err", err.Error()))
		userInfo = &client.UserInfo{
			UserID: userID.String(),
			Email:  "unknown@email.com",
			Name:   "Unknown User",
		}
	}

	// Получаем информацию о владельце вишлиста
	ownerInfo, err := s.authClient.GetUserInfo(ctx, wl.OwnerID.String())
	if err != nil {
		log.Error("Failed to get user info", slog.String("err", err.Error()))
		ownerInfo = &client.UserInfo{
			UserID: wl.OwnerID.String(),
			Email:  "unknown@email.com",
			Name:   "Unknown Owner",
		}
	}

	// Получаем информацию о том, кто ИЗНАЧАЛЬНО бронировал (для уведомлений)
	bookedByInfo, err := s.authClient.GetUserInfo(ctx, lastOwnerID.String())
	if err != nil {
		log.Error("Failed to get user info", slog.String("err", err.Error()))
		bookedByInfo = &client.UserInfo{
			UserID: lastOwnerID.String(),
			Email:  "unknown@email.com",
			Name:   "Unknown User",
		}
	}

	// Определяем причину отмены
	reason := "cancelled_by_user"
	if userInfo.UserID == ownerInfo.UserID {
		reason = "cancelled_by_owner"
	}

	event := producer.ItemUnbookedEvent{
		ItemID:       item.ID.String(),
		ItemName:     item.Name,
		WishlistID:   wl.ID.String(),
		WishlistName: wl.Title,

		UnbookedBy:      userInfo.UserID,
		UnbookedByEmail: userInfo.Email,
		UnbookedByName:  userInfo.Name,

		OwnerID:    ownerInfo.UserID,
		OwnerEmail: ownerInfo.Email,
		OwnerName:  ownerInfo.Name,

		BookedBy:      bookedByInfo.UserID,
		BookedByEmail: bookedByInfo.Email,
		BookedByName:  bookedByInfo.Name,

		Reason: reason,
	}

	go func() {
		if err := s.kafkaProd.PublishItemUnbooked(context.Background(), event); err != nil {
			log.Error("Failed to publish unbooked event", slog.String("err", err.Error()))
		}
	}()
}
