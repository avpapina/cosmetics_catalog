package database

import (
	"cosmetics_catalog/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	// Подключение к SQLite (файл будет создан автоматически)
	db, err := gorm.Open(sqlite.Open("cosmetics.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db

	// Автомиграция - создание таблиц
	err = db.AutoMigrate(
		&models.Brand{},
		&models.Category{},
		&models.Subcategory{},
		&models.Product{},
	)
	if err != nil {
		return err
	}

	return nil
}
