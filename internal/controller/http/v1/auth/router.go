package auth

import (
	"github.com/Markard/wordka/internal/usecase"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func GetRouter(
	logger logger.Interface,
	val *validator.Validate,
	authUC *usecase.Auth,
) *chi.Mux {
	r := chi.NewRouter()
	c := NewController(authUC, logger, val)

	r.Post("/register", c.Register)
	r.Post("/login", c.Login)

	return r
}
