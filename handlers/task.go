package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"todo-api/models"
)

type TaskHandler struct {
	DB *gorm.DB
}

func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

// ✅ CreateTask: создаёт задачу для текущего авторизованного пользователя
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task

	// Получаем user_id из JWT (установлен в middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	// Привязываем задачу к пользователю
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.UserID = userID.(uint)

	if err := h.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать задачу"})
		return
	}

	task.UserID = userID.(uint)
	task.SetTitle("[Создано] " + task.Title)

	c.JSON(http.StatusCreated, task)
}

// ✅ GetTasks: возвращает только задачи текущего пользователя
func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	var tasks []models.Task
	if err := h.DB.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить задачи"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// Остальные методы пока без авторизации (можно доработать позже)
func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	var task models.Task
	if err := h.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	var task models.Task
	if err := h.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}

	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.Title = input.Title
	task.Description = input.Description
	task.Completed = input.Completed

	if err := h.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить задачу"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	if err := h.DB.Delete(&models.Task{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить задачу"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Задача удалена"})
}

func (h *TaskHandler) DebugSliceUsage(c *gin.Context) {
	// Имитируем получение задач из базы
	tasks := []models.Task{
		{Title: "Задача 1"},
		{Title: "Задача 2"},
	}

	// Добавим задачу динамически
	newTask := models.Task{Title: "Задача 3"}
	tasks = append(tasks, newTask)

	// Вернём клиенту
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) DebugMapUsage(c *gin.Context) {
	tasks := []models.Task{
		{Model: gorm.Model{ID: 1}, Title: "Задача 1"},
		{Model: gorm.Model{ID: 2}, Title: "Задача 2"},
		{Model: gorm.Model{ID: 3}, Title: "Задача 3"},
	}

	taskMap := make(map[uint]models.Task)
	for _, task := range tasks {
		taskMap[task.ID] = task
	}

	task := taskMap[2]

	c.JSON(http.StatusOK, gin.H{
		"task_2":  task,
		"taskMap": taskMap,
	})
}
