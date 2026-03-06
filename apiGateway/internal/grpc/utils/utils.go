package utils

import (
	"fmt"
	"gateway/internal/domain/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcToDomainError(grpcError error) error {
	st, ok := status.FromError(grpcError)
	if !ok {
		return grpcError
	}

	// Получаем код и сообщение
	code := st.Code()
	message := st.Message()
	// Обрабатываем конкретные коды
	switch code {
	case codes.NotFound:
		return fmt.Errorf("%s: %w", message, types.ErrorNotFound)
	case codes.InvalidArgument:
		return fmt.Errorf("%s: %w", message, types.ErrorInvalidArgument)
	case codes.PermissionDenied:
		return fmt.Errorf("%s: %w", message, types.ErrorAccessDenied)
	case codes.AlreadyExists:
		return fmt.Errorf("%s: %w", message, types.ErrorAlreadyExist)
	default: // codes.Internal
		return fmt.Errorf("%s: %w", message, types.ErrorInternal)
	}
}
