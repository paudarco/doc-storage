//go:build !windows
// +build !windows

package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/paudarco/doc-storage/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg config.Server, handler http.Handler) error {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	s.httpServer = &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 20, // 1 MB
		Handler:        handler,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
