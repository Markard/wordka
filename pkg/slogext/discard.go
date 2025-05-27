package slogext

import (
	"context"
	"log/slog"
)

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Enabled(ctx context.Context, rec slog.Level) bool {
	return true
}

func (h *DiscardHandler) Handle(ctx context.Context, rec slog.Record) error {
	return nil
}

func (h *DiscardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) WithGroup(name string) slog.Handler {
	return &DiscardHandler{}
}
