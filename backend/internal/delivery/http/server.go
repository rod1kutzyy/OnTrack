package http

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rod1kutzyy/OnTrack/internal/delivery/http/handlers"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string, todoHandler *handlers.TodoHandler) *Server {
	router := gin.Default()
	RegisterTodoRoutes(router, todoHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	return &Server{httpServer: srv}
}

func (s *Server) Run() error {
	log.Printf("Server started on address %s\n", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
