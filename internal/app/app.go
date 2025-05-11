package app

import (
	"fmt"
	"github.com/Markard/wordka/config"
	"github.com/Markard/wordka/internal/controller/http"
	"github.com/Markard/wordka/internal/infra/middleware"
	"github.com/Markard/wordka/internal/infra/middleware/jwt"
	serviceJwt "github.com/Markard/wordka/internal/infra/service/jwt"
	"github.com/Markard/wordka/internal/repo"
	"github.com/Markard/wordka/internal/usecase"
	"github.com/Markard/wordka/internal/usecase/auth"
	"github.com/Markard/wordka/pkg/httpserver"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/Markard/wordka/pkg/postgres"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func Run(setup *config.Setup) {
	// Initialize logger
	logFile, err := os.OpenFile(setup.Config.Log.FilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
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
	lgr := logger.New(setup.Config.Log.Level, setup.Config.Log.CallerSkipFrameCount, logFile)

	// Validator
	val, err := httpserver.NewValidator()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("could not initiate validator")
	}

	// Repository PostgreSQL
	db := postgres.New(setup.Env.PgDSN, lgr.ZerologLogger())
	defer func() {
		err := db.Close()
		if err != nil {
			lgr.Error(err)
		}
	}()
	authRepo := repo.NewAuthRepository(db)
	gameRepo := repo.NewGameRepository(db)

	// Use cases
	jwtService := serviceJwt.NewService(setup.Env.ES256PrivateKey, setup.Env.ES256PublicKey)
	useCases := &usecase.UseCases{
		AuthUseCase: auth.NewAuth(authRepo, jwtService),
		GameUseCase: usecase.NewGameUseCase(gameRepo),
	}

	// Middleware
	middlewares := &middleware.Middlewares{
		JwtAuthenticator: jwt.Authenticator(jwtService, authRepo),
	}

	// HTTP Server
	httpServer := httpserver.New(setup.Config.HttpServer.Address, setup.Config.HttpServer.IdleTimeout)
	http.SetupRouter(httpServer.Router, setup, lgr, val, middlewares, useCases)

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
