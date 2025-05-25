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
	"github.com/Markard/wordka/pkg/logger"
	"github.com/Markard/wordka/pkg/postgres"
	"os"
	"os/signal"
	"syscall"
)

func Run(setup *config.Setup) {
	lgr := logger.New(setup.Config.Log.Level)

	// Validator
	val, err := validator.NewValidator()
	if err != nil {
		lgr.Fatal(fmt.Errorf("could not initiate validator: %w", err))
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
		GameUseCase: game.NewGameUseCase(gameRepo),
	}

	// Middleware
	middlewares := &middleware.Middlewares{
		JwtAuthenticator: jwt.Authenticator(jwtService, authRepo),
	}

	// HTTP Server
	httpServer := server.New(setup.Config.HttpServer.Address, setup.Config.HttpServer.IdleTimeout)
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
