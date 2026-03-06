package models

import (
	"time"

	"github.com/google/uuid"
)

type Wishlist struct {
	ID          uuid.UUID
	OwnerID     uuid.UUID
	Title       string
	Description string
	IsPublic    bool
	Items       []Item
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewWishList(ownerID uuid.UUID, title, description string, isPublic bool) Wishlist {
	now := time.Now()
	return Wishlist{
		ID:          uuid.New(),
		OwnerID:     ownerID,
		Title:       title,
		Description: description,
		IsPublic:    isPublic,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (w *Wishlist) Update(title, description string, isPublic bool) {
	w.Title = title
	w.Description = description
	w.IsPublic = isPublic
	w.UpdatedAt = time.Now()
}

func (w *Wishlist) CanView(userID uuid.UUID) bool {
	if w.IsPublic {
		return true
	} else if userID == w.OwnerID { // Только владелец может глянуть
		return true
	}

	return false
}

func (w *Wishlist) CanEdit(userID uuid.UUID) bool {
	if userID == w.OwnerID { // Только владелец может изменять
		return true
	}

	return false
}
