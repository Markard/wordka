package slogext

import (
	"context"
	"log/slog"
)

type logCtx struct {
	UserId int64
}

const logCtxKey = "logCtx"

type ContextEnricher struct {
	next slog.Handler
}

func NewContextEnricher(next slog.Handler) slog.Handler {
	return &ContextEnricher{next: next}
}

func (h *ContextEnricher) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *ContextEnricher) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(logCtxKey).(logCtx); ok {
		rec.Add("userId", c.UserId)
	}
	return h.next.Handle(ctx, rec)
}

func (h *ContextEnricher) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextEnricher{next: h.next.WithAttrs(attrs)}
}

func (h *ContextEnricher) WithGroup(name string) slog.Handler {
	return &ContextEnricher{next: h.next.WithGroup(name)}
}

func WithLogUserID(ctx context.Context, userId int64) context.Context {
	if c, ok := ctx.Value(logCtxKey).(logCtx); ok {
		c.UserId = userId
		return context.WithValue(ctx, logCtxKey, c)
	}
	return context.WithValue(ctx, logCtxKey, logCtx{UserId: userId})
}
