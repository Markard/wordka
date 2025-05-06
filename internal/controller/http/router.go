package http

import (
	"github.com/Markard/wordka/config"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewRouter(rootRouter *chi.Mux, cfg *config.Config, logger logger.Interface) {
	rootRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
}
