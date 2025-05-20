package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type Server struct {
	Router     *chi.Mux
	context    context.Context
	httpServer *http.Server
	notify     chan error
}

func New(address string, idleTimeout time.Duration) *Server {
	router := chi.NewRouter()

	httpServer := &http.Server{
		Addr:        address,
		Handler:     router,
		IdleTimeout: idleTimeout,
	}

	return &Server{
		Router:     router,
		httpServer: httpServer,
		context:    context.Background(),
		notify:     make(chan error, 1),
	}
}

func (server *Server) Start() {
	go func() {
		server.notify <- server.httpServer.ListenAndServe()
		close(server.notify)
	}()
}

func (server *Server) Notify() <-chan error {
	return server.notify
}

func (server *Server) Shutdown() error {
	return server.httpServer.Shutdown(server.context)
}
