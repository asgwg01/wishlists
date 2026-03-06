package utils

import (
	"log/slog"
	"wishlistService/internal/domain/models"
)

func WishlistToSlog(wl models.Wishlist) slog.Attr {
	result := slog.Group(
		"wishlist",
		slog.String("id", wl.ID.String()),
		slog.String("owner_id", wl.OwnerID.String()),
		slog.String("title", wl.Title),
		slog.String("description", wl.Description),
		slog.Bool("is_public", wl.IsPublic),
		slog.Time("create_at", wl.CreatedAt),
		slog.Time("update_at", wl.UpdatedAt),
		slog.Int("items_count", len(wl.Items)),
	)
	return result
}

func ItemToSlog(item models.Item) slog.Attr {
	result := slog.Group(
		"item",
		slog.String("id", item.ID.String()),
		slog.String("wishlist_id", item.WishlistID.String()),
		slog.String("name", item.Name),
		slog.String("description", item.Description),
		slog.String("image_url", item.ImageURL),
		slog.String("product_url", item.ProductURL),
		slog.String("price", item.Price.String()),
		slog.Any("booked_by", item.BookedBy),
		slog.Any("booked_at", item.BookedAt),
		slog.Time("create_at", item.CreatedAt),
		slog.Time("update_at", item.UpdatedAt),
	)
	return result
}
