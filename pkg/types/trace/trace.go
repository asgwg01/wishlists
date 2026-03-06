package trace

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"google.golang.org/grpc/metadata"
)

const (
	// HeaderKey - ключ для trace ID в заголовках
	HeaderKey = "x-trace-id"

	// MetadataKey - ключ для trace ID в gRPC metadata
	MetadataKey = "trace-id"
)

// GenerateTraceID генерирует новый trace ID
func GenerateTraceID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetTraceID извлекает trace ID из контекста
func GetTraceID(ctx context.Context) string {
	// Пробуем получить из gRPC metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(MetadataKey); len(values) > 0 {
			return values[0]
		}
	}

	// Пробуем получить из context values
	if traceID, ok := ctx.Value(MetadataKey).(string); ok {
		return traceID
	}

	return ""
}

// WithTraceID добавляет trace ID в контекст
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, MetadataKey, traceID)
}

// InjectIntoGRPC добавляет trace ID в gRPC metadata для исходящих запросов
func InjectIntoGRPC(ctx context.Context) context.Context {
	traceID := GetTraceID(ctx)
	if traceID == "" {
		traceID = GenerateTraceID()
	}

	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		// Копируем существующую metadata
		md = md.Copy()
	} else {
		md = metadata.New(make(map[string]string))
	}

	md.Set(MetadataKey, traceID)

	return metadata.NewOutgoingContext(ctx, md)
}

// // InjectIntoHTTP добавляет trace ID в HTTP заголовки
// func InjectIntoHTTP(ctx context.Context, header map[string][]string) {
// 	traceID := GetTraceID(ctx)
// 	if traceID != "" {
// 		header[HeaderKey] = []string{traceID}
// 	}
// }
