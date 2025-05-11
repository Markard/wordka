package v1

import (
	"github.com/Markard/wordka/internal/controller/http/v1/auth"
	"github.com/Markard/wordka/internal/controller/http/v1/game"
	"github.com/Markard/wordka/internal/usecase"
	projectMiddleware "github.com/Markard/wordka/pkg/httpserver/middleware"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func CreateRouter(
	logger logger.Interface,
	val *validator.Validate,
	tokenVerifier projectMiddleware.TokenVerifier,
	authUC *usecase.Auth,
) *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/", auth.CreateRouter(logger, val, authUC))
	r.Group(func(r chi.Router) {
		r.Mount("/games", game.CreateRouter(logger, val))
	}).With(projectMiddleware.JwtVerifier(tokenVerifier))

	return r
}
