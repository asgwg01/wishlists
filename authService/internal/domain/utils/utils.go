package utils

import (
	"authService/internal/domain/models"
	"log/slog"
)

func UserToSlog(user models.User) slog.Attr {
	result := slog.Group(
		"user",
		slog.String("id", user.ID.String()),
		slog.String("email", user.Email),
		slog.String("name", user.Name),
		//slog.String("pwd", string(user.PasswordHash)),
		slog.Time("create_at", user.CreateAt),
		slog.Time("update_at", user.UpdateAt),
	)
	return result
}
