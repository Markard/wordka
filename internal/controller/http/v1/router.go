package v1

import (
	"github.com/Markard/wordka/internal/controller/http/v1/auth"
	"github.com/Markard/wordka/internal/controller/http/v1/game"
	"github.com/Markard/wordka/internal/infra/middleware"
	"github.com/Markard/wordka/internal/usecase"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func CreateRouter(
	logger logger.Interface,
	val *validator.Validate,
	middlewares *middleware.Middlewares,
	useCases *usecase.UseCases,
) *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/", auth.CreateRouter(logger, val, useCases.AuthUseCase))
	r.Group(func(r chi.Router) {
		r.Use(middlewares.JwtAuthenticator)

		r.Mount("/games", game.CreateRouter(logger, val))
	})

	return r
}
