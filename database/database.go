package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
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

func ConnectPGX() (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:1111@localhost:5432/todo?sslmode=disable"
	}
	return pgxpool.New(context.Background(), dsn)
}
