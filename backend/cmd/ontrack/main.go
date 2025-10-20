package main

import (
	"log"

	"github.com/rod1kutzyy/OnTrack/internal/config"
	"github.com/rod1kutzyy/OnTrack/internal/delivery/http"
	"github.com/rod1kutzyy/OnTrack/internal/delivery/http/handlers"
	"github.com/rod1kutzyy/OnTrack/internal/infrastructure/database"
	"github.com/rod1kutzyy/OnTrack/internal/repository"
	"github.com/rod1kutzyy/OnTrack/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("error getting database instance: %v", err)
	}
	defer sqlDB.Close()

	if cfg.App.Env == "development" {
		err := database.RunMigration(db)
		if err != nil {
			log.Fatalf("unable to run migrations: %v", err)
		}

		log.Println("Database migrated successfully")
	}

	todoRepo := repository.NewTodoRepository(db)
	todoUseCase := usecase.NewTodoUsecase(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoUseCase)
	server := http.NewServer(cfg.App.Port, todoHandler)

	if err := server.Run(); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
