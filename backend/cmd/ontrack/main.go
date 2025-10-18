package main

import (
	"log"

	"github.com/rod1kutzyy/OnTrack/internal/config"
	"github.com/rod1kutzyy/OnTrack/internal/infrastructure/database"
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

	if cfg.App.Env == "development" {
		err := database.RunMigration(db)
		if err != nil {
			log.Fatalf("unable to run migrations: %v", err)
		}

		log.Println("Database migrated successfully")
	}
}
