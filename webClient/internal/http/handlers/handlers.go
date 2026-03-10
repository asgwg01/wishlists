package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorDTO struct {
	Error string `json:"error"`
}

type RegisterRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthDTO struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	UserID       string `json:"user_id"`
	ExpiresAt    int64  `json:"expires_at"`
}

type UserInfoDTO struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type CreateWishlistRequestDTO struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type UpdateWishlistRequestDTO = CreateWishlistRequestDTO

type WishlistDTO struct {
	ID          string `json:"id"`
	OwnerID     string `json:"owner_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type WishlistListDTO struct {
	Wishlists   []WishlistDTO `json:"wishlists"`
	TotalCount  int32         `json:"total_count"`
	TotalPages  int32         `json:"total_pages"`
	CurrentPage int32         `json:"current_page"`
}

type CreateItemRequestDTO struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	ProductURL  string `json:"product_url"`
	Price       int64  `json:"price"`
}

type UpdateItemRequestDTO = CreateItemRequestDTO

type ItemDTO struct {
	ID          string  `json:"id"`
	WishlistID  string  `json:"wishlist_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	ProductURL  string  `json:"product_url"`
	Price       int64   `json:"price"`
	BookedBy    *string `json:"booked_by,omitempty"`
	BookedAt    *int64  `json:"booked_at,omitempty"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}

type ItemListDTO struct {
	Items       []ItemDTO `json:"items"`
	TotalCount  int32     `json:"total_count"`
	TotalPages  int32     `json:"total_pages"`
	CurrentPage int32     `json:"current_page"`
}

type BookingDTO struct {
	ItemID     string `json:"item_id"`
	WishlistID string `json:"wishlist_id"`
	UserID     string `json:"user_id"`
	BookedAt   int64  `json:"booked_at"`
	ItemName   string `json:"item_name,omitempty"`
}

type BookingListDTO struct {
	Bookings   []BookingDTO `json:"bookings"`
	TotalCount int32        `json:"total_count"`
}

func SendError(log *slog.Logger, w http.ResponseWriter, status int, writeMsg, logMsg string) {
	log.Error(logMsg)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	errDto := ErrorDTO{
		Error: writeMsg,
	}
	if err := json.NewEncoder(w).Encode(errDto); err != nil {
		log.Error("Error Encode DTO", slog.String("err", err.Error()))
		return
	}
}
