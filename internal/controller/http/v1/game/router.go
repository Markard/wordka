package game

import (
	"github.com/Markard/wordka/internal/usecase/game"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func CreateRouter(
	logger logger.Interface,
	val *validator.Validate,
	useCase *game.UseCase,
) *chi.Mux {
	r := chi.NewRouter()
	c := NewController(useCase, logger, val)

	r.Get("/current", c.GetCurrentGame)
	r.Post("/current", c.CreateGame)

	return r
}
