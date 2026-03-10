package clients

import "httpClient/internal/http/handlers"

type IWishlistGatewayService interface {
	// Auth api calls
	Register(req handlers.RegisterRequestDTO) (handlers.AuthDTO, error)
	Login(req handlers.LoginRequestDTO) (handlers.AuthDTO, error)
	GetCurrentUser(token string) (handlers.UserInfoDTO, error)
	GetUserInfo(token string, userID string) (handlers.UserInfoDTO, error)

	// Wishlists api calls
	CreateWishlist(token string, req handlers.CreateWishlistRequestDTO) (handlers.WishlistDTO, error)
	GetUserWishlists(token, userID string, page, limit int32) (handlers.WishlistListDTO, error)
	GetPublicWishlists(token string, page, limit int32) (handlers.WishlistListDTO, error)
	GetWishlist(token, wishlistID string) (handlers.WishlistDTO, error)
	UpdateWishlist(token, wishlistID string, req handlers.UpdateWishlistRequestDTO) (handlers.WishlistDTO, error)
	DeleteWishlist(token, wishlistID string) error

	// Item api calls
	CreateItem(token, wishlistID string, req handlers.CreateItemRequestDTO) (handlers.ItemDTO, error)
	GetItems(token, wishlistID string, page, limit int32) (handlers.ItemListDTO, error)
	GetItem(token, itemID string) (handlers.ItemDTO, error)
	UpdateItem(token, itemID string, req handlers.UpdateItemRequestDTO) (handlers.ItemDTO, error)
	DeleteItem(token, itemID string) error

	BookItem(token, itemID string) (handlers.ItemDTO, error)
	UnbookItem(token, itemID string) (handlers.ItemDTO, error)

	// Booking api calls
	GetUserBookings(token string) (handlers.BookingListDTO, error)
}
