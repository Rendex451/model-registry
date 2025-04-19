package config

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var storage *gorm.DB

func InitDB(databasePath string, models ...any) error {
	dir := filepath.Dir(databasePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory for database: %v", err)
		return err
	}

	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to SQLite database: %v", err)
		return err
	}

	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
		return err
	}

	storage = db
	log.Println("SQLite database connection established")

	return nil
}

func GetDB() *gorm.DB {
	return storage
}
