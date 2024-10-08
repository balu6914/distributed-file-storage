package db

import (
	"distributed-file-storage/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("files.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto-migrate the FilePart model to create the necessary table
	DB.AutoMigrate(&models.FilePart{})
	return nil
}
