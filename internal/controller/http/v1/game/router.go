package game

import (
	"github.com/Markard/wordka/internal/usecase/game"
	"github.com/Markard/wordka/pkg/http/validator"
	"github.com/go-chi/chi/v5"
)

func CreateRouter(val validator.ProjectValidator, useCase *game.UseCase) *chi.Mux {
	r := chi.NewRouter()
	c := NewController(useCase, val)

	r.Get("/", c.GetCurrentGame)
	r.Post("/", c.CreateGame)
	r.Post("/guess", c.Guess)

	return r
}
