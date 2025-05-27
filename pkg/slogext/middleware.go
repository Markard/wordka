package slogext

import (
	"log/slog"
)

type NewHandlerMiddleware func(next slog.Handler) slog.Handler
