package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/pkgerrors"
	"net/http"
	"os"
	"strings"
	"time"
)

type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(err error)
	Fatal(err error)
	RequestLogger(next http.Handler) http.Handler
	ZerologLogger() *zerolog.Logger
}

type Logger struct {
	logger *zerolog.Logger
}

var _ Interface = (*Logger)(nil)

func New(level string, callerSkipFrameCount int) *Logger {
	setGlobalLevel(level)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + callerSkipFrameCount).
		Logger()

	return &Logger{
		logger: &logger,
	}
}

func setGlobalLevel(level string) {
	switch strings.ToLower(level) {
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func (logger *Logger) Debug(message interface{}, args ...interface{}) {
	logger.msg("debug", message, args...)
}

func (logger *Logger) Info(message string, args ...interface{}) {
	logger.log(message, args...)
}

func (logger *Logger) Warn(message string, args ...interface{}) {
	logger.log(message, args...)
}

func (logger *Logger) Error(err error) {
	logger.logger.Error().Stack().Err(err).Msg("")
}

func (logger *Logger) Fatal(err error) {
	logger.logger.
		Fatal().
		Err(err).
		Msg("")

	os.Exit(1)
}

func (logger *Logger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case string:
		logger.log(msg, args...)
	default:
		logger.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}

func (logger *Logger) log(message string, args ...interface{}) {
	if len(args) == 0 {
		logger.logger.Info().Msg(message)
	} else {
		logger.logger.Info().Msgf(message, args...)
	}
}

func (logger *Logger) RequestLogger(next http.Handler) http.Handler {
	h := hlog.NewHandler(*logger.logger)

	accessHandler := hlog.AccessHandler(
		func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status_code", status).
				Int("response_size_bytes", size).
				Dur("elapsed_ms", duration).
				Msg("incoming request")
		},
	)

	userAgentHandler := hlog.UserAgentHandler("http_user_agent")
	remoteAddrHandler := hlog.RemoteAddrHandler("ip")
	refererHandler := hlog.RefererHandler("referer")
	requestIDHandler := hlog.RequestIDHandler("req_id", "Request-Id")

	return h(
		accessHandler(
			userAgentHandler(
				remoteAddrHandler(
					refererHandler(
						requestIDHandler(
							next,
						),
					),
				),
			),
		),
	)
}

func (logger *Logger) ZerologLogger() *zerolog.Logger {
	return logger.logger
}
