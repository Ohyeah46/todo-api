package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"todo-api/models"
)

func ConnectPostgres() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=1111 dbname=todo port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Task{})
}
