package internal

import (
	"github.com/gorilla/mux"
)

type RouteInjector interface {
	Register(r *mux.Router, logger *Logger)
}

type Server struct {
	Router *mux.Router
	Logger *Logger
}

func NewServer(logger *Logger, router *mux.Router, routes []RouteInjector) *Server {
	server := &Server{
		Router: router,
		Logger: logger,
	}

	for _, r := range routes {
		r.Register(server.Router, server.Logger)
	}

	return server
}
