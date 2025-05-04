package http

import (
	"github.com/Markard/wordka/config"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func SetupRouter(router *chi.Mux, cfg *config.Config, logger logger.Interface) {
	router.Use(middleware.Timeout(cfg.HttpServer.Timeout))
	router.Use(middleware.Recoverer)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
}
