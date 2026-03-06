package stubstorage

import (
	"log/slog"
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

// metods storage.IStorage
