package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// ErrorDTO DTO для описания ошибки
type ErrorDTO struct {
	// Error содержит строку с расшифровкой того что пошло не так
	// minLength: 1
	// maxLength: 100
	// example: item not found
	Error string `json:"error"`
}

// RegisterRequestDTO DTO для запроса регистрации нового пользователя
type RegisterRequestDTO struct {
	// Email пользователя
	// required: true
	// format: email
	// example: user@example.com
	Email string `json:"email" validate:"required,email"`
	// Пароль пользователя (минимум 4 символа)
	// required: true
	// minLength: 4
	// example: password123
	Password string `json:"password" validate:"required,min=6"`
	// Имя, ник пользователя
	// required: true
	// example: Иван
	Name string `json:"name" validate:"required"`
}

// LoginRequestDTO DTO для запроса входа пользователя
type LoginRequestDTO struct {
	// Email пользователя
	// required: true
	// format: email
	// example: user@example.com
	Email string `json:"email" validate:"required,email"`
	// Пароль пользователя (минимум 4 символа)
	// required: true
	// minLength: 4
	// example: password123
	Password string `json:"password" validate:"required"`
}

// AuthDTO DTO для описания данных аутентификации
type AuthDTO struct {
	// JWT токен доступа
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Token string `json:"token"`
	// Токен для обновления доступа
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	RefreshToken string `json:"refresh_token,omitempty"`
	// UserID UUID пользователя
	// format: uuid
	// example: 123e4567-e89b-12d3-a456-426614174000
	UserID string `json:"user_id"`
	// ExpiresAt Время окончания токена (Unix timestamp)
	// format: date-time
	// example: 1700000000
	ExpiresAt int64 `json:"expires_at"`
}

// UserInfoDTO DTO для описания данных пользователя
type UserInfoDTO struct {
	// UserID UUID пользователя
	// format: uuid
	// example: 123e4567-e89b-12d3-a456-426614174000
	UserID string `json:"user_id"`
	// Email пользователя
	// required: true
	// format: email
	// example: user@example.com
	Email string `json:"email"`
	// Name Имя, ник пользователя
	// required: true
	// example: Иван
	Name string `json:"name"`
}

// CreateWishlistRequestDTO DTO для запроса создания нового вишлиста
type CreateWishlistRequestDTO struct {
	// Title Название вишлиста
	// required: true
	// minLength: 1
	// maxLength: 80
	// example: На день рождения
	Title string `json:"title" validate:"required"`
	// Description Описание вишлиста
	// maxLength: 1000
	// example: Подарки на мой день рождения
	Description string `json:"description"`
	// IsPublic Флаг публичности вишлиста
	// example: true
	IsPublic bool `json:"is_public"`
}

// UpdateWishlistRequestDTO DTO для запроса обновления вишлиста
type UpdateWishlistRequestDTO = CreateWishlistRequestDTO

// WishlistDTO DTO для описания вишлиста
type WishlistDTO struct {
	// ID UUID вишлиста
	// format: uuid
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID string `json:"id"`
	// OwnerID UUID пользователя - владельца вишлиста
	// format: uuid
	// example: 123e4567-e89b-12d3-a456-426614174000
	OwnerID string `json:"owner_id"`
	// Title Название вишлиста
	// required: true
	// minLength: 1
	// maxLength: 80
	// example: На день рождения
	Title string `json:"title"`
	// Description Описание вишлиста
	// maxLength: 1000
	// example: Подарки на мой день рождения
	Description string `json:"description"`
	// IsPublic Флаг публичности вишлиста
	// example: true
	IsPublic bool `json:"is_public"`
	// CreatedAt Дата и время создания (Unix timestamp)
	// format: date-time
	// example: 1700000000
	CreatedAt int64 `json:"created_at"`
	// UpdatedAt Дата и время последнего обновления (Unix timestamp)
	// format: date-time
	// example: 1700000000
	UpdatedAt int64 `json:"updated_at"`
}

// WishlistListDTO DTO для описания списка вишлистов с пагинацией
type WishlistListDTO struct {
	// Wishlists Список вишлистов
	Wishlists []WishlistDTO `json:"wishlists"`
	// TotalCount Общее количество вишлистов в списке Wishlists
	// example: 100
	TotalCount int32 `json:"total_count"`
	// TotalPages Общее количество страниц
	// example: 10
	TotalPages int32 `json:"total_pages"`
	// CurrentPage Текущая страница
	// example: 1
	CurrentPage int32 `json:"current_page"`
}

// CreateItemRequestDTO DTO для запроса создания нового айтема
type CreateItemRequestDTO struct {
	// Name Название айтема
	// required: true
	// minLength: 1
	// maxLength: 80
	// example: Сертификат на метание топоров
	Name string `json:"name" validate:"required"`
	// Description Описание айтема
	// maxLength: 1000
	// example: Сертификат, желательно на четверых, и чтобы красивый был и с топорами
	Description string `json:"description"`
	// ImageURL URL изображения
	// format: uri
	// example: https://example.com/image.jpg
	ImageURL string `json:"image_url"`
	// ProductURL URL на то где можно приобрести
	// format: uri
	// example: https://example.com/product/123
	ProductURL string `json:"product_url"`
	// Price Цена в копейках (или предполагаемая цена)
	// minimum: 0
	// example: 50000
	Price int64 `json:"price"`
}

// CreateItemRequestDTO DTO для запроса обновления айтема
type UpdateItemRequestDTO = CreateItemRequestDTO

// ItemDTO DTO для описания айтема, пожелания, элемента вишлиста
type ItemDTO struct {
	// ID UUID айтема
	// format: uuid
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID string `json:"id"`
	// ID UUID вишлиста в котором находится айтем
	// format: uuid
	// example: 123e4567-e89b-12d3-a456-426614174000
	WishlistID string `json:"wishlist_id"`
	// Name Название айтема
	// required: true
	// minLength: 1
	// maxLength: 80
	// example: Сертификат на метание топоров
	Name string `json:"name"`
	// Description Описание айтема
	// maxLength: 1000
	// example: Сертификат, желательно на четверых, и чтобы красивый был и с топорами
	Description string `json:"description"`
	// ImageURL URL изображения
	// format: uri
	// example: https://example.com/image.jpg
	ImageURL string `json:"image_url"`
	// ProductURL URL на то где можно приобрести
	// format: uri
	// example: https://example.com/product/123
	ProductURL string `json:"product_url"`
	// Price Цена в копейках (или предполагаемая цена)
	// minimum: 0
	// example: 50000
	Price int64 `json:"price"`
	// BookedBy UUID пользователя, забронировавшего предмет, если забронирован
	// format: uuid
	// nullable: true
	// example: 123e4567-e89b-12d3-a456-426614174000
	BookedBy *string `json:"booked_by,omitempty"`
	// BookedAt Дата и время бронирования, если есть
	// format: date-time
	// nullable: true
	// example: 1700000000
	BookedAt *int64 `json:"booked_at,omitempty"`
	// CreatedAt Дата и время создания (Unix timestamp)
	// format: date-time
	// example: 1700000000
	CreatedAt int64 `json:"created_at"`
	// UpdatedAt Дата и время последнего обновления (Unix timestamp)
	// format: date-time
	// example: 1700000000
	UpdatedAt int64 `json:"updated_at"`
}

// ItemListDTO DTO для описания списка айтемов с пагинацией
type ItemListDTO struct {
	// Items Список айтемов
	Items []ItemDTO `json:"items"`
	// TotalCount Общее количество айтемов в списке WishlItemsists
	// example: 100
	TotalCount int32 `json:"total_count"`
	// TotalPages Общее количество страниц
	// example: 10
	TotalPages int32 `json:"total_pages"`
	// CurrentPage Текущая страница
	// example: 1
	CurrentPage int32 `json:"current_page"`
}

// / BookingDTO DTO для описания бронирования айтема пользователем
type BookingDTO struct {
	// ItemID UUID айтема
	// format: uuid
	// example: 123e4567-e89b-12d3-a456-426614174000
	ItemID string `json:"item_id"`
	// WishlistID UUID вишлиста в котором находитя айтем
	// format: uuid
	// example: 123e4567-e89b-12d3-a456-426614174000
	WishlistID string `json:"wishlist_id"`
	// UserID UUID пользователя который забронировал айтем
	// format: uuid
	// example: 123e4567-e89b-12d3-a456-426614174000
	UserID string `json:"user_id"`
	// BookedAt Дата и время бронирования
	// format: date-time
	// example: 1700000000
	BookedAt int64 `json:"booked_at"`
	// ItemName Название айтема, доп инфо чтобы не делать лишний запрос позднее
	// required: true
	// minLength: 1
	// maxLength: 80
	// example: Сертификат на метание топоров
	ItemName string `json:"item_name,omitempty"`
}

// ItemListDTO DTO для описания списка бронирований с пагинацией
type BookingListDTO struct {
	// Bookings Список бронирований
	Bookings []BookingDTO `json:"bookings"`
	// TotalCount Общее количество бронирований в списке Bookings
	// example: 100
	TotalCount int32 `json:"total_count"`
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
