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
	"github.com/Markard/wordka/internal/usecase/game"
	"github.com/Markard/wordka/pkg/http/server"
	"github.com/Markard/wordka/pkg/http/validator"
	"github.com/Markard/wordka/pkg/postgres"
	"github.com/Markard/wordka/pkg/slogext"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Run(setup *config.Setup) {
	logger := slogext.SetupLogger(setup.Env.AppEnv)
	slog.SetDefault(logger)

	// Validator
	val, err := validator.NewValidator()
	if err != nil {
		slogext.Fatal(logger, err)
	}

	// Repository PostgreSQL
	db := postgres.New(setup.Env.PgDSN, logger)
	defer func() {
		if err := db.Close(); err != nil {
			slogext.Fatal(logger, err)
		}
	}()
	authRepo := repo.NewAuthRepository(db)
	gameRepo := repo.NewGameRepository(db)

	// Use cases
	jwtService := serviceJwt.NewService(setup.Env.ES256PrivateKey, setup.Env.ES256PublicKey)
	useCases := &usecase.UseCases{
		AuthUseCase: auth.NewAuth(authRepo, jwtService),
		GameUseCase: game.NewGameUseCase(gameRepo),
	}

	// Middleware
	middlewares := &middleware.Middlewares{
		JwtAuthenticator: jwt.Authenticator(jwtService, authRepo, logger),
	}

	// HTTP Server
	httpServer := server.New(setup.Config.HttpServer.Address, setup.Config.HttpServer.IdleTimeout)
	http.SetupRouter(httpServer.Router, setup, val, middlewares, useCases)

	// Start Http Server
	httpServer.Start()
	logger.Info("Wordka:Start", "address", setup.Config.HttpServer.Address, "env", setup.Env.AppEnv)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("Wordka:Signal", "signal", s.String())
	case err = <-httpServer.Notify():
		slogext.Error(logger, fmt.Errorf("Wordka:Running | Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		slogext.Error(logger, fmt.Errorf("Wordka:Shutdown | Error: %w", err))
	}
	logger.Info("Wordka:Shutdown")
}
