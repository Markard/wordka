package auth

import (
	"github.com/Markard/wordka/internal/usecase/auth"
	"github.com/Markard/wordka/pkg/httpserver"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/chi/v5"
)

func CreateRouter(
	logger logger.Interface,
	val httpserver.ProjectValidator,
	useCase *auth.UseCase,
) *chi.Mux {
	r := chi.NewRouter()
	c := NewController(useCase, logger, val)

	r.Post("/register", c.Register)
	r.Post("/login", c.Login)

	return r
}
