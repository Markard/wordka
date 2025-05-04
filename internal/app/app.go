package app

import (
	"github.com/Markard/wordka/config"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/rs/zerolog/log"
	"os"
)

func Run(env *config.Env, cfg *config.Config) {
	logFile, err := os.OpenFile(cfg.Log.FilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("could not open log file")
	}
	defer func(logFile *os.File) {
		if err := logFile.Close(); err != nil {
			log.With().Err(err)
		}
	}(logFile)
	lgr := logger.New(cfg.Log.Level, cfg.Log.CallerSkipFrameCount, logFile)
	lgr.Info("starting wordka app")
}
