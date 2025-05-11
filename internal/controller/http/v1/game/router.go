package game

import (
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func CreateRouter(
	logger logger.Interface,
	val *validator.Validate,
) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, "Hello games!")
	})

	return r
}
