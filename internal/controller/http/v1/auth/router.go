package auth

import (
	"github.com/Markard/wordka/internal/usecase"
	"github.com/Markard/wordka/pkg/httpserver"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/chi/v5"
)

func CreateRouter(
	logger logger.Interface,
	val httpserver.ProjectValidator,
	authUseCase *usecase.Auth,
) *chi.Mux {
	r := chi.NewRouter()
	c := NewController(authUseCase, logger, val)

	r.Post("/register", c.Register)
	r.Post("/login", c.Login)

	return r
}
