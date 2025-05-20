package http

import (
	"github.com/Markard/wordka/config"
	"github.com/Markard/wordka/internal/controller/http/v1"
	projectMiddleware "github.com/Markard/wordka/internal/infra/middleware"
	"github.com/Markard/wordka/internal/usecase"
	"github.com/Markard/wordka/pkg/http/validator"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"net/http"
)

func SetupRouter(
	router *chi.Mux,
	setup *config.Setup,
	logger logger.Interface,
	val validator.ProjectValidator,
	middlewares *projectMiddleware.Middlewares,
	useCases *usecase.UseCases,
) {
	router.Use(logger.RequestLogger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(setup.Config.HttpServer.Timeout))

	router.Get("/robots.txt", robotsTxt)
	router.Mount("/v1", v1.CreateRouter(logger, val, middlewares, useCases))
}

func robotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	robotTxt := `User-agent: *
Disallow: /`
	render.PlainText(w, r, robotTxt)
}
