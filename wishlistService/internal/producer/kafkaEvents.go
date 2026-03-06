package producer

import (
	"time"

	"github.com/google/uuid"
)

// EventType - тип события
type EventType string

const (
	EventTypeItemBooked   EventType = "item.booked"
	EventTypeItemUnbooked EventType = "item.unbooked"
)

type BaseEvent struct {
	EventID      string    `json:"event_id"`
	EventType    EventType `json:"event_type"`
	EventVersion string    `json:"event_version"`
	Timestamp    int64     `json:"timestamp"`
}

// ItemBookedEvent - событие бронирования
type ItemBookedEvent struct {
	BaseEvent
	ItemID       string `json:"item_id"`
	ItemName     string `json:"item_name"`
	WishlistID   string `json:"wishlist_id"`
	WishlistName string `json:"wishlist_name"`

	// Кто забронировал
	BookedBy      string `json:"booked_by"`
	BookedByEmail string `json:"booked_by_email"`
	BookedByName  string `json:"booked_by_name"`

	// Владелец вишлиста
	OwnerID    string `json:"owner_id"`
	OwnerEmail string `json:"owner_email"`
	OwnerName  string `json:"owner_name"`
}

// ItemUnbookedEvent - событие отмены
type ItemUnbookedEvent struct {
	BaseEvent
	ItemID       string `json:"item_id"`
	ItemName     string `json:"item_name"`
	WishlistID   string `json:"wishlist_id"`
	WishlistName string `json:"wishlist_name"`

	// Кто отменил
	UnbookedBy      string `json:"unbooked_by"`
	UnbookedByEmail string `json:"unbooked_by_email"`
	UnbookedByName  string `json:"unbooked_by_name"`

	// Владелец вишлиста
	OwnerID    string `json:"owner_id"`
	OwnerEmail string `json:"owner_email"`
	OwnerName  string `json:"owner_name"`

	// Кто изначально бронировал (для уведомлений при отмене)
	BookedBy      string `json:"booked_by"`
	BookedByEmail string `json:"booked_by_email"`
	BookedByName  string `json:"booked_by_name"`

	Reason string `json:"reason"`
}

func NewItemBookedEvent(itemID, itemName, wishlistID, wishlistName,
	bookedBy, bookedByEmail, bookedByName,
	ownerID, ownerEmail, ownerName string) ItemBookedEvent {
	return ItemBookedEvent{
		BaseEvent: BaseEvent{
			EventID:      uuid.New().String(),
			EventType:    EventTypeItemBooked,
			EventVersion: "1.0",
			Timestamp:    time.Now().Unix(),
		},
		ItemID:        itemID,
		ItemName:      itemName,
		WishlistID:    wishlistID,
		WishlistName:  wishlistName,
		BookedBy:      bookedBy,
		BookedByEmail: bookedByEmail,
		BookedByName:  bookedByName,
		OwnerID:       ownerID,
		OwnerEmail:    ownerEmail,
		OwnerName:     ownerName,
	}
}

func NewItemUnbookedEvent(itemID, itemName, wishlistID, wishlistName,
	unbookedBy, unbookedByEmail, unbookedByName,
	ownerID, ownerEmail, ownerName,
	bookedBy, bookedByEmail, bookedByName,
	reason string) ItemUnbookedEvent {
	return ItemUnbookedEvent{
		BaseEvent: BaseEvent{
			EventID:      uuid.New().String(),
			EventType:    EventTypeItemUnbooked,
			EventVersion: "1.0",
			Timestamp:    time.Now().Unix(),
		},
		ItemID:          itemID,
		ItemName:        itemName,
		WishlistID:      wishlistID,
		WishlistName:    wishlistName,
		UnbookedBy:      unbookedBy,
		UnbookedByEmail: unbookedByEmail,
		UnbookedByName:  unbookedByName,
		OwnerID:         ownerID,
		OwnerEmail:      ownerEmail,
		OwnerName:       ownerName,
		BookedBy:        bookedBy,
		BookedByEmail:   bookedByEmail,
		BookedByName:    bookedByName,
		Reason:          reason,
	}
}
