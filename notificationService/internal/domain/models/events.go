package models

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
