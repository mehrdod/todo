package todo

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Startup(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ": " + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   40 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	s.httpServer.Shutdown(ctx)
}
