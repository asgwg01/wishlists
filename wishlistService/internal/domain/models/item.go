package models

import (
	"time"
	"wishlistService/internal/domain/types"

	"github.com/asgwg01/wishlists/pkg/types/price"

	"github.com/google/uuid"
)

type Item struct {
	ID          uuid.UUID
	WishlistID  uuid.UUID
	Name        string
	Description string
	ImageURL    string
	ProductURL  string
	Price       price.Price
	BookedBy    *uuid.UUID // может быть пустым
	BookedAt    *time.Time // может быть пустым
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewItem(wishlistID uuid.UUID, name, description, imageURL, productURL string, price price.Price) Item {
	now := time.Now()
	return Item{
		ID:          uuid.New(),
		WishlistID:  wishlistID,
		Name:        name,
		Description: description,
		ImageURL:    imageURL,
		ProductURL:  productURL,
		Price:       price,
		BookedBy:    nil,
		BookedAt:    nil,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (i *Item) Update(name, description, imageURL, productURL string, price price.Price) {
	i.Name = name
	i.Description = description
	i.ImageURL = imageURL
	i.ProductURL = productURL
	i.Price = price
	i.UpdatedAt = time.Now()
}

func (i *Item) IsBooked() (bool, *uuid.UUID) {
	if i.BookedBy != nil {
		return true, i.BookedBy
	}

	return false, nil
}

func (i *Item) Book(userID uuid.UUID) error {
	if i.BookedBy != nil {
		return types.ErrorItemAlreadyBooked
	}

	now := time.Now()
	i.BookedBy = &userID
	i.BookedAt = &now
	i.UpdatedAt = now

	return nil
}

func (i *Item) Unbook(userID uuid.UUID) error {
	if i.BookedBy == nil {
		return types.ErrorItemNotBooked
	}

	if *i.BookedBy != userID {
		return types.ErrorAccessDenied
	}

	i.BookedBy = nil
	i.BookedAt = nil
	i.UpdatedAt = time.Now()

	return nil
}
