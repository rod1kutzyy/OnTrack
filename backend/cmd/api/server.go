package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rod1kutzyy/OnTrack/internal/config"
	"github.com/rod1kutzyy/OnTrack/internal/logger"
)

type Server struct {
	httpServer *http.Server
	config     *config.Config
}

func NewServer(cfg *config.Config, router *gin.Engine) *Server {
	httpServer := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &Server{
		httpServer: httpServer,
		config:     cfg,
	}
}

func (s *Server) Start() <-chan error {
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)
		logger.Logger.Infof("Starting server on %s", s.httpServer.Addr)

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to start server: %w", err)
		}
	}()

	return errChan
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Logger.Info("Shutting down server gracefully...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logger.Logger.Info("Server stopped gracefully")
	return nil
}

func (s *Server) WaitForShutdownSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Logger.Infof("Received shutdown signal: %v", sig)
}

func (s *Server) GracefulShutdown(cleanup func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		logger.Logger.Errorf("Error during server shutdown: %v", err)
	}

	if cleanup != nil {
		if err := cleanup(); err != nil {
			logger.Logger.Errorf("Error during cleanup: %v", err)
			return err
		}
	}

	logger.Logger.Info("Application stopped successfully")
	return nil
}
