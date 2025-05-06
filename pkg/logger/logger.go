package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Logger struct {
	logger *zerolog.Logger
}

var _ Interface = (*Logger)(nil)

func New(level string, callerSkipFrameCount int, logFile *os.File) *Logger {
	setGlobalLevel(level)

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

func (logger *Logger) Error(message interface{}, args ...interface{}) {
	if logger.logger.GetLevel() == zerolog.DebugLevel {
		logger.Debug(message, args...)
	}

	logger.msg("error", message, args...)
}

func (logger *Logger) Fatal(message interface{}, args ...interface{}) {
	logger.msg("fatal", message, args...)

	os.Exit(1)
}

func (logger *Logger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		logger.log(msg.Error(), args...)
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
