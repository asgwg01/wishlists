package logger

import (
	"context"
	"log/slog"
)

// MultiHandler для объединения нескольких handlers slog
type MultiHandler struct {
	Handlers []slog.Handler
}

func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.Handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, handler := range h.Handlers {
		if err := handler.Handle(ctx, record); err != nil {
			return err
		}
	}
	return nil
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.Handlers))
	for i, handler := range h.Handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return &MultiHandler{Handlers: newHandlers}
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.Handlers))
	for i, handler := range h.Handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &MultiHandler{Handlers: newHandlers}
}
