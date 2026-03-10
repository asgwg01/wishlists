package apihandlers

import (
	"httpClient/internal/http/clients"
	"log/slog"
)

type WishlistGatewayHandlers struct {
	AuthHandlers     *AuthHandlers
	WishlistHandlers *WishlistHandlers
	ItemHandlers     *ItemHandlers
	BookingHandlers  *BookingHandlers
}

func NewWishlistGatewayHandlers(log *slog.Logger, wlgClient clients.IWishlistGatewayService) *WishlistGatewayHandlers {
	return &WishlistGatewayHandlers{
		AuthHandlers:     NewAuthHandlers(log, wlgClient),
		WishlistHandlers: NewWishlistHandlers(log, wlgClient),
		ItemHandlers:     NewItemHandlers(log, wlgClient),
		BookingHandlers:  NewBookingHandlers(log, wlgClient),
	}
}
