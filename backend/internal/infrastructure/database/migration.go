package database

import (
	"github.com/rod1kutzyy/OnTrack/internal/entity"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) error {
	err := db.AutoMigrate(&entity.Todo{})
	if err != nil {
		return err
	}

	return nil
}
