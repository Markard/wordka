package httpserver

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

type Server struct {
	Router          *chi.Mux
	HttpServer      *http.Server
	Context         context.Context
	notify          chan error
	shutdownTimeout time.Duration
}

func New(
	address string,
	timeout time.Duration,
	idleTimeout time.Duration,
	shutdownTimeout time.Duration,
) *Server {
	router := chi.NewRouter()
	router.Use(middleware.Timeout(timeout))

	httpServer := &http.Server{
		Addr:        address,
		Handler:     router,
		IdleTimeout: idleTimeout,
	}

	return &Server{
		Router:          router,
		HttpServer:      httpServer,
		Context:         context.Background(),
		notify:          make(chan error, 1),
		shutdownTimeout: shutdownTimeout,
	}
}

func (server *Server) Start() {
	go func() {
		server.notify <- server.HttpServer.ListenAndServe()
		close(server.notify)
	}()
}

func (server *Server) Notify() <-chan error {
	return server.notify
}

func (server *Server) Shutdown() error {
	return server.HttpServer.Shutdown(server.Context)
}
