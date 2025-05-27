package slogext

import (
	"context"
	"github.com/Markard/wordka/config/env"
	"log/slog"
	"os"
	"runtime"
	"time"
)

func SetupLogger(appEnv string) *slog.Logger {
	var handler slog.Handler

	switch appEnv {
	case env.Dev:
		handler = NewPrettyHandler(&slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}, NewContextEnricher)
	case env.Test:
		handler = NewDiscardHandler()
	case env.Prod:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

func Fatal(logger *slog.Logger, err error) {
	Error(logger, err)
	os.Exit(1)
}

func Error(logger *slog.Logger, err error) {
	if !logger.Enabled(context.Background(), slog.LevelError) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelError, err.Error(), pcs[0])
	_ = logger.Handler().Handle(context.Background(), r)
}
