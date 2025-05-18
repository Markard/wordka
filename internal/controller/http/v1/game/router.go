package game

import (
	"github.com/Markard/wordka/internal/usecase/game"
	"github.com/Markard/wordka/pkg/httpserver"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/chi/v5"
)

func CreateRouter(
	logger logger.Interface,
	val httpserver.ProjectValidator,
	useCase *game.UseCase,
) *chi.Mux {
	r := chi.NewRouter()
	c := NewController(useCase, logger, val)

	r.Get("/", c.GetCurrentGame)
	r.Post("/", c.CreateGame)
	r.Post("/guess", c.Guess)

	return r
}
