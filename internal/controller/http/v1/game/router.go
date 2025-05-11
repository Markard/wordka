package game

import (
	"github.com/Markard/wordka/internal/entity"
	"github.com/Markard/wordka/internal/infra/middleware/jwt"
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
		currentUser, _ := r.Context().Value(jwt.CurrentUserCtxKey).(*entity.User)

		render.JSON(w, r, currentUser)

		// Получи из контекста user id
		// Вызови юзкейс с user id
		render.JSON(w, r, "Hello games!")
	})

	return r
}
