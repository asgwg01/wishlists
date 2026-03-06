package wishlistService

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"wishlistService/internal/domain/models"
	"wishlistService/internal/domain/types"
	"wishlistService/internal/storage"

	"github.com/asgwg01/wishlists/pkg/types/trace"

	"github.com/google/uuid"
)

type WshlistService struct {
	log             *slog.Logger
	wishlistStorage storage.IWishlistStorage
}

type IWshlistService interface {
	Create(ctx context.Context, ownerID uuid.UUID, title, description string, isPublic bool) (models.Wishlist, error)
	Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (models.Wishlist, error)
	GetUserWishlists(ctx context.Context, ownerID uuid.UUID, requestingUserID uuid.UUID) ([]models.Wishlist, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, title, description string, isPublic bool) (models.Wishlist, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	PublicWishlists(ctx context.Context, page, pageSize int) ([]models.Wishlist, int, error)
}

func New(
	log *slog.Logger,
	wishlistStorage storage.IWishlistStorage,
) *WshlistService {
	return &WshlistService{
		log:             log,
		wishlistStorage: wishlistStorage,
	}
}

func (s *WshlistService) Create(ctx context.Context,
	ownerID uuid.UUID,
	title string,
	description string,
	isPublic bool,
) (models.Wishlist, error) {
	const logPrefix = "service.wishlistService.Create"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Create wishlist")
	if title == "" {
		return models.Wishlist{}, fmt.Errorf("title is required")
	}

	wishlist := models.NewWishList(ownerID, title, description, isPublic)

	wishlist, err := s.wishlistStorage.CreateWishlist(ctx, wishlist)
	if err != nil {
		if errors.Is(err, storage.ErrorWishlistAlreadyExist) {
			log.Error("Wishlist already exist", slog.String("err", err.Error()))
			return models.Wishlist{}, types.ErrorWishlistAlreadyExist
		}
		log.Error("Error create wishlist", slog.String("err", err.Error()))
		return models.Wishlist{}, fmt.Errorf("failed to create wishlist: %w", err)
	}

	return wishlist, nil
}
func (s *WshlistService) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (models.Wishlist, error) {
	const logPrefix = "service.wishlistService.Get"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Get wishlist", slog.String("uuid", id.String()))

	wl, err := s.wishlistStorage.GetWishlistByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrorWishlistNotExist) {
			return models.Wishlist{}, types.ErrorWishlistNotFound
		}
		log.Error("Error get wishlist", slog.String("err", err.Error()))
		return models.Wishlist{}, err
	}

	if !wl.CanView(userID) {
		return models.Wishlist{}, types.ErrorAccessDenied
	}

	return wl, nil
}
func (s *WshlistService) GetUserWishlists(ctx context.Context,
	ownerID uuid.UUID,
	requestingUserID uuid.UUID,
) ([]models.Wishlist, error) {
	const logPrefix = "service.wishlistService.Get"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Get wishlists", slog.String("ownerID", ownerID.String()))

	includePrivate := false
	if ownerID == requestingUserID {
		includePrivate = true
	}

	wl, err := s.wishlistStorage.GetWishlistsByOwnerID(ctx, ownerID, includePrivate)
	if err != nil {
		log.Error("Error get wishlist", slog.String("err", err.Error()))
		return []models.Wishlist{}, err
	}

	return wl, nil
}
func (s *WshlistService) Update(ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
	title string,
	description string,
	isPublic bool,
) (models.Wishlist, error) {
	const logPrefix = "service.wishlistService.Update"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Update wishlist", slog.String("id", id.String()))

	if title == "" {
		log.Error("Error title can not be empty")
		return models.Wishlist{}, fmt.Errorf("Error title can not be empty")
	}

	wl, err := s.wishlistStorage.GetWishlistByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrorWishlistNotExist) {
			return models.Wishlist{}, types.ErrorWishlistNotFound
		}
		return models.Wishlist{}, fmt.Errorf("failed to get wishlist: %w", err)
	}

	if !wl.CanEdit(userID) {
		return models.Wishlist{}, types.ErrorAccessDenied
	}

	wl.Update(title, description, isPublic)
	wl, err = s.wishlistStorage.UpdateWishlist(ctx, wl)
	if err != nil {
		if errors.Is(err, storage.ErrorWishlistNotExist) {
			return models.Wishlist{}, types.ErrorWishlistNotFound
		}
		return models.Wishlist{}, fmt.Errorf("failed to update wishlist: %w", err)
	}

	return wl, nil
}
func (s *WshlistService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	const logPrefix = "service.wishlistService.Delete"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Delete wishlist", slog.String("id", id.String()))

	wl, err := s.wishlistStorage.GetWishlistByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrorWishlistNotExist) {
			return types.ErrorWishlistNotFound
		}
		return fmt.Errorf("failed to get wishlist: %w", err)
	}

	if !wl.CanEdit(userID) {
		return types.ErrorAccessDenied
	}

	err = s.wishlistStorage.DeleteWishlist(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrorWishlistNotExist) {
			return types.ErrorWishlistNotFound
		}
		return fmt.Errorf("failed to delete wishlist: %w", err)
	}

	return nil
}
func (s *WshlistService) PublicWishlists(ctx context.Context, page, pageSize int) ([]models.Wishlist, int, error) {
	const logPrefix = "service.wishlistService.ListPublic"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Get public wishlists", slog.Int("page", page), slog.Int("pageSize", pageSize))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	wishlists, total, err := s.wishlistStorage.ListPublicWishlists(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list public wishlists: %w", err)
	}

	return wishlists, total, nil

}
