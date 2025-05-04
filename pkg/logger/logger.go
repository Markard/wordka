package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"strings"
)

type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(err error)
	Fatal(err error)
}

type Logger struct {
	logger *zerolog.Logger
}

var _ Interface = (*Logger)(nil)

func New(level string, callerSkipFrameCount int, logFile *os.File) *Logger {
	setGlobalLevel(level)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	logger := zerolog.New(multi).
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
