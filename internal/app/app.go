package app

import (
	"fmt"
	"github.com/Markard/wordka/config"
	"github.com/Markard/wordka/internal/controller/http"
	"github.com/Markard/wordka/pkg/httpserver"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/Markard/wordka/pkg/postgres"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func Run(env *config.Env, cfg *config.Config) {
	// Initialize logger
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

	// Repository PostgreSQL
	db := postgres.New(env.PgDSN, lgr.ZerologLogger())
	defer func() {
		err := db.Close()
		if err != nil {
			lgr.Error(err)
		}
	}()

	// HTTP Server
	httpServer := httpserver.New(cfg.HttpServer.Address, cfg.HttpServer.IdleTimeout)
	http.SetupRouter(httpServer.Router, cfg, lgr)

	// Start Http Server
	lgr.Info("app - Run - httpServer.Start")
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		lgr.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		lgr.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		lgr.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
