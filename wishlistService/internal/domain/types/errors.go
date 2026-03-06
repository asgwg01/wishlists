package types

import "errors"

var (
	ErrorItemNotFound      = errors.New("item not found")
	ErrorItemAlreadyExist  = errors.New("item already exist")
	ErrorItemAlreadyBooked = errors.New("item already booked")
	ErrorItemNotBooked     = errors.New("item not booked")

	ErrorWishlistNotFound     = errors.New("wishlist not found")
	ErrorWishlistAlreadyExist = errors.New("wishlist already exist")

	ErrorAccessDenied = errors.New("access denied")
)
