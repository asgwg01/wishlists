package stubstorage

import (
	"context"
	"log/slog"
	"time"
	"wishlistService/internal/domain/models"
	"wishlistService/internal/domain/utils"

	"github.com/google/uuid"
)

type Storage struct {
	log *slog.Logger
}

func NewStorage(log *slog.Logger) (*Storage, error) {
	l := log.With(
		slog.String("STUB!", "stubStorage"),
	)
	return &Storage{log: l}, nil
}

// IWishlistStorage
func (s *Storage) CreateWishlist(ctx context.Context, wishlist models.Wishlist) (models.Wishlist, error) {
	const logPrefix = "stubstorage.Storage.CreateWishlist"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", utils.WishlistToSlog(wishlist))

	return models.Wishlist{}, nil
}
func (s *Storage) GetWishlistByID(ctx context.Context, id uuid.UUID) (models.Wishlist, error) {
	const logPrefix = "stubstorage.Storage.GetWishlistByID"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", slog.String("uuid", id.String()))

	return models.Wishlist{}, nil
}
func (s *Storage) GetWishlistsByOwnerID(ctx context.Context, ownerID uuid.UUID, includePrivate bool) ([]models.Wishlist, error) {
	const logPrefix = "stubstorage.Storage.GetWishlistsByOwnerID"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", slog.String("owner uuid", ownerID.String()), slog.Bool("includePrivate", includePrivate))

	return []models.Wishlist{}, nil
}
func (s *Storage) UpdateWishlist(ctx context.Context, wishlist models.Wishlist) (models.Wishlist, error) {
	const logPrefix = "stubstorage.Storage.UpdateWishlist"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", utils.WishlistToSlog(wishlist))

	return models.Wishlist{}, nil
}
func (s *Storage) DeleteWishlist(ctx context.Context, id uuid.UUID) error {
	const logPrefix = "stubstorage.Storage.DeleteWishlist"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", slog.String("uuid", id.String()))

	return nil
}
func (s *Storage) ListPublicWishlists(ctx context.Context, limit, offset int) ([]models.Wishlist, int, error) {
	const logPrefix = "stubstorage.Storage.ListPublicWishlists"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", slog.Int("limit", limit), slog.Int("offset", offset))

	return []models.Wishlist{}, 42, nil
}

// IItemStorage
func (s *Storage) CreateItem(ctx context.Context, item models.Item) (models.Item, error) {
	const logPrefix = "stubstorage.Storage.CreateItem"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", utils.ItemToSlog(item))

	return models.Item{}, nil
}
func (s *Storage) GetItemByID(ctx context.Context, id uuid.UUID) (models.Item, error) {
	const logPrefix = "stubstorage.Storage.GetItemByID"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", slog.String("uuid", id.String()))

	return models.Item{}, nil
}
func (s *Storage) GetItemsByWishlistID(ctx context.Context, wishlistID uuid.UUID, limit, offset int) ([]models.Item, int, error) {
	const logPrefix = "stubstorage.Storage.GetItemsByWishlistID"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", slog.String("wl uuid", wishlistID.String()), slog.Int("wl limit", limit), slog.Int("wl offset", offset))

	return []models.Item{}, 42, nil
}
func (s *Storage) UpdateItem(ctx context.Context, item models.Item) (models.Item, error) {
	const logPrefix = "stubstorage.Storage.UpdateItem"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", utils.ItemToSlog(item))

	return models.Item{}, nil
}
func (s *Storage) DeleteItem(ctx context.Context, id uuid.UUID) error {
	const logPrefix = "stubstorage.Storage.DeleteItem"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", slog.String("uuid", id.String()))

	return nil
}
func (s *Storage) UpdateItemBooking(ctx context.Context, itemID uuid.UUID, bookedBy *uuid.UUID, bookedAt *time.Time) (models.Item, error) {
	const logPrefix = "stubstorage.Storage.UpdateItemBooking"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", slog.String("itemID", itemID.String()), slog.String("bookedBy", bookedBy.String()))

	return models.Item{}, nil
}

func (s *Storage) GetBookedItemsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Item, error) {
	const logPrefix = "stubstorage.Storage.GetBookedItemsByUser"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", slog.String("itemID", userID.String()))

	return []models.Item{}, nil
}
