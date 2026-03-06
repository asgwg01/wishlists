package storage

import (
	"authService/internal/domain/models"
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrorUserNotExist = errors.New("user is not exist")
	ErrorEmailExist   = errors.New("email already exists")
)

type IStorage interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type ICreator interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
}

type IGeter interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type IUpdater interface {
	UpdateUser(ctx context.Context, user models.User) (models.User, error)
}

type IDeleter interface {
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
