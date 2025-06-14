package internal

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type RouteInjector interface {
	Register(r *mux.Router, logger *Logger)
}

type Server struct {
	http.Server
	Router *mux.Router
	Logger *Logger
}

func NewServer(logger *Logger, router *mux.Router, routes []RouteInjector) *Server {
	server := &Server{
		Router: router,
		Logger: logger,
	}
	server.Handler = server.Router

	for _, r := range routes {
		r.Register(server.Router, server.Logger)
	}

	return server
}

func (s *Server) Start(addr string) error {
	s.Addr = addr

	s.Logger.Info("Starting server", zap.String("address", addr))
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.Logger.Error("Failed to start server", zap.String("error", err.Error()))
		return err
	}

	return nil
}
