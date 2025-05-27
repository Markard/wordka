package auth

import (
	"github.com/Markard/wordka/internal/usecase/auth"
	"github.com/Markard/wordka/pkg/http/validator"
	"github.com/go-chi/chi/v5"
)

func CreateRouter(
	val validator.ProjectValidator,
	useCase *auth.UseCase,
) *chi.Mux {
	r := chi.NewRouter()
	c := NewController(useCase, val)

	r.Post("/register", c.Register)
	r.Post("/login", c.Login)

	return r
}
