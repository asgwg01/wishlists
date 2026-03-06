package stubstorage

import (
	"authService/internal/domain/models"
	"authService/internal/domain/utils"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type Storage struct {
	log *slog.Logger
}

func NewStorage(log *slog.Logger) (*Storage, error) {
	l := log.With(
		slog.String("STUB!", "stubStorage"),
	)
	return &Storage{log: l}, nil
}
func (s *Storage) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	const logPrefix = "stubstorage.Storage.CreateUser"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", utils.UserToSlog(user))
	return user, nil
}
func (s *Storage) GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	const logPrefix = "stubstorage.Storage.GetUserByID"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", slog.String("user uuid", id.String()))
	return models.User{}, nil
}
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const logPrefix = "stubstorage.Storage.GetUserByEmail"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", slog.String("user email", email))
	return models.User{}, nil
}
func (s *Storage) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	const logPrefix = "stubstorage.Storage.UpdateUser"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", slog.String("user uuid", user.ID.String()))
	return user, nil
}
func (s *Storage) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const logPrefix = "stubstorage.Storage.DeleteUser"
	log := s.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Call STUB", slog.String("user uuid", id.String()))
	return nil
}

// metods storage.IStorage
