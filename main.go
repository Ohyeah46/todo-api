package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"todo-api/database"
	"todo-api/handlers"
)

func main() {
	// Подключение к GORM
	db, err := database.ConnectPostgres()
	if err != nil {
		log.Fatalf("Ошибка подключения к GORM: %v", err)
	}
	database.Migrate(db)

	// Подключение к pgxpool для sqlc
	pool, err := database.ConnectPGX()
	if err != nil {
		log.Fatalf("Ошибка подключения к PGX: %v", err)
	}
	defer pool.Close()

	handlers.InitTaskWorker()

	router := gin.Default()

	authHandler := handlers.NewAuthHandler(db)
	taskHandler := handlers.NewTaskHandler(db)
	taskSQLCHandler := handlers.NewTaskSQLCHandler(pool)

	// Публичные роуты
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	// Отладочные
	router.GET("/debug/slice", taskHandler.DebugSliceUsage)
	router.GET("/debug/map", taskHandler.DebugMapUsage)
	router.GET("/async-example", taskHandler.AsyncProcessExample)

	auth := router.Group("/")
	auth.Use(handlers.AuthMiddleware())

	auth.POST("/enqueue", taskHandler.EnqueueTask)
	auth.POST("/tasks", taskHandler.CreateTask) // GORM
	auth.GET("/tasks", taskHandler.GetTasks)    // GORM
	auth.GET("/tasks/:id", taskHandler.GetTask) // GORM
	auth.PUT("/tasks/:id", taskHandler.UpdateTask)
	auth.DELETE("/tasks/:id", taskHandler.DeleteTask)

	// SQLC маршруты (тестовые)
	auth.GET("/sqlc/tasks", taskSQLCHandler.GetTasks)
	auth.GET("/sqlc/tasks/:id", taskSQLCHandler.GetTask)
	auth.POST("/sqlc/tasks", taskSQLCHandler.CreateTask)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
