package storage

import (
	"context"
	"errors"
	"time"
	"wishlistService/internal/domain/models"

	"github.com/google/uuid"
)

var (
	ErrorWishlistNotExist     = errors.New("wishlist is not exist")
	ErrorItemNotExist         = errors.New("item is not exist")
	ErrorWishlistAlreadyExist = errors.New("wishlist already exist")
	ErrorItemAlreadyExist     = errors.New("item already exist")
)

type IWishlistStorage interface {
	CreateWishlist(ctx context.Context, wishlist models.Wishlist) (models.Wishlist, error)
	GetWishlistByID(ctx context.Context, id uuid.UUID) (models.Wishlist, error)
	GetWishlistsByOwnerID(ctx context.Context, ownerID uuid.UUID, includePrivate bool) ([]models.Wishlist, error)
	UpdateWishlist(ctx context.Context, wishlist models.Wishlist) (models.Wishlist, error)
	DeleteWishlist(ctx context.Context, id uuid.UUID) error
	ListPublicWishlists(ctx context.Context, limit, offset int) ([]models.Wishlist, int, error)
}

type IItemStorage interface {
	CreateItem(ctx context.Context, item models.Item) (models.Item, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (models.Item, error)
	GetItemsByWishlistID(ctx context.Context, wishlistID uuid.UUID, limit, offset int) ([]models.Item, int, error)
	UpdateItem(ctx context.Context, item models.Item) (models.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
	UpdateItemBooking(ctx context.Context, itemID uuid.UUID, bookedBy *uuid.UUID, bookedAt *time.Time) (models.Item, error)
	GetBookedItemsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Item, error)
}
