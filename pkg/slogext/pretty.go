package slogext

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
)

const (
	timeFormat = "[15:04:05.000]"

	reset = "\033[0m"

	cyan        = 36
	darkGray    = 90
	lightRed    = 91
	lightGreen  = 92
	lightYellow = 93
	lightCyan   = 96
)

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

type PrettyHandler struct {
	h slog.Handler
	b *bytes.Buffer
	m *sync.Mutex
}

func NewPrettyHandler(opts *slog.HandlerOptions, hm NewHandlerMiddleware) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	b := &bytes.Buffer{}
	return &PrettyHandler{
		h: hm(
			slog.NewJSONHandler(b, &slog.HandlerOptions{
				Level:       opts.Level,
				AddSource:   opts.AddSource,
				ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
			}),
		),
		b: b,
		m: &sync.Mutex{},
	}
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{h: h.h.WithGroup(name), b: h.b, m: h.m}
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = colorize(darkGray, level)
	case slog.LevelInfo:
		level = colorize(cyan, level)
	case slog.LevelWarn:
		level = colorize(lightYellow, level)
	case slog.LevelError:
		level = colorize(lightRed, level)
	}

	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("error when marshaling attrs: %w", err)
	}

	fmt.Println(
		colorize(lightGreen, r.Time.Format(timeFormat)),
		level,
		colorize(darkGray, r.Message),
		colorize(lightCyan, string(b)),
	)

	return nil
}

func (h *PrettyHandler) computeAttrs(ctx context.Context, r slog.Record) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()
	if err := h.h.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	if err := json.Unmarshal(h.b.Bytes(), &attrs); err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}
	return attrs, nil
}

func suppressDefaults(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}
