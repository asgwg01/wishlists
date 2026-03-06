package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	Name         string
	PasswordHash []byte
	CreateAt     time.Time
	UpdateAt     time.Time
}
