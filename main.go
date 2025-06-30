package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"todo-api/database"
	"todo-api/handlers"
)

func main() {
	// Подключение к БД
	db, err := database.ConnectPostgres()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	// Миграция моделей User и Task
	database.Migrate(db)

	// Инициализация роутера Gin
	router := gin.Default()

	// Инициализация хендлеров
	authHandler := handlers.NewAuthHandler(db)
	taskHandler := handlers.NewTaskHandler(db)

	// Публичные маршруты
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.GET("/debug/slice", taskHandler.DebugSliceUsage)
	router.GET("/debug/map", taskHandler.DebugMapUsage)

	// Группа защищённых маршрутов
	auth := router.Group("/")
	auth.Use(handlers.AuthMiddleware())

	auth.POST("/tasks", taskHandler.CreateTask)
	auth.GET("/tasks", taskHandler.GetTasks)
	auth.GET("/tasks/:id", taskHandler.GetTask)
	auth.PUT("/tasks/:id", taskHandler.UpdateTask)
	auth.DELETE("/tasks/:id", taskHandler.DeleteTask)

	// Запуск сервера
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
