package http

import (
	"github.com/Markard/wordka/config"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"net/http"
)

func SetupRouter(router *chi.Mux, cfg *config.Config, logger logger.Interface) {
	router.Use(logger.RequestLogger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(cfg.HttpServer.Timeout))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	router.Get("/robots.txt", robotsTxt)
}

func robotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	robotTxt := `User-agent: *
Disallow: /`
	render.PlainText(w, r, robotTxt)
}
