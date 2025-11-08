package main

import (
	"fmt"
	"os"

	_ "github.com/rod1kutzyy/OnTrack/docs"
	"github.com/rod1kutzyy/OnTrack/internal/config"
	"github.com/rod1kutzyy/OnTrack/internal/domain"
	"github.com/rod1kutzyy/OnTrack/internal/handler"
	"github.com/rod1kutzyy/OnTrack/internal/infrastructure/database"
	"github.com/rod1kutzyy/OnTrack/internal/logger"
	"github.com/rod1kutzyy/OnTrack/internal/repository/postgres"
	"github.com/rod1kutzyy/OnTrack/internal/usecase"
	"github.com/rod1kutzyy/OnTrack/internal/validator"
)

// @title OnTrack API
// @version 1.0
// @description The REST API of the OnTrack service is an application for managing tasks.
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	if err := logger.Init(cfg.Logger.Level); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	logger.Logger.Info("=== Starting Application ===")
	logger.Logger.Infof("Environment: %s", cfg.Logger.Level)

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.Logger.Fatalf("Failed to initialize database: %v", err)
	}

	if err := db.AutoMigrate(&domain.Todo{}); err != nil {
		logger.Logger.Fatalf("Failed to run database migrations: %v", err)
	}

	todoRepo := postgres.NewTodoRepository(db.GetDB())
	todoUseCase := usecase.NewTodoUseCase(todoRepo)
	todoValidator := validator.NewTodoValidator()
	todoHandler := handler.NewTodoHandler(todoUseCase, todoValidator)

	router := SetupRouter(cfg, todoHandler)
	srv := NewServer(cfg, router)

	errChan := srv.Start()

	go func() {
		if err := <-errChan; err != nil {
			logger.Logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	srv.WaitForShutdownSignal()

	cleanup := func() error {
		return db.Close()
	}

	if err := srv.GracefulShutdown(cleanup); err != nil {
		logger.Logger.Errorf("Error during shutdown: %v", err)
		os.Exit(1)
	}

	logger.Logger.Info("=== Application Exited Successfully ===")
}
