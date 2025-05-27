package v1

import (
	"github.com/Markard/wordka/internal/controller/http/v1/auth"
	"github.com/Markard/wordka/internal/controller/http/v1/game"
	"github.com/Markard/wordka/internal/infra/middleware"
	"github.com/Markard/wordka/internal/usecase"
	"github.com/Markard/wordka/pkg/http/validator"
	"github.com/go-chi/chi/v5"
)

func CreateRouter(
	val validator.ProjectValidator,
	middlewares *middleware.Middlewares,
	useCases *usecase.UseCases,
) *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/", auth.CreateRouter(val, useCases.AuthUseCase))
	r.Group(func(r chi.Router) {
		r.Use(middlewares.JwtAuthenticator)

		r.Mount("/games/current", game.CreateRouter(val, useCases.GameUseCase))
	})

	return r
}
