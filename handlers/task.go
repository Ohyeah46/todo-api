package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/jackc/pgx/v5/pgxpool" // pgxpool для sqlc подключения
	"todo-api/internal/db"            // sqlc сгенерированный пакет
	"todo-api/models"
)

// ---------- Хендлер с GORM ----------

type TaskHandler struct {
	DB *gorm.DB
}

func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.UserID = userID.(uint)
	task.SetTitle("[Создано] " + task.Title)

	if err := h.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать задачу"})
		return
	}

	go func(t models.Task) {
		log.Printf("Task created asynchronously: ID=%d, Title=%s\n", t.ID, t.Title)
	}(task)

	c.JSON(http.StatusCreated, task)
}

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

	summary := task.ShortSummary()
	overdue := task.IsOverdue()

	c.JSON(http.StatusOK, gin.H{
		"task":    task,
		"summary": summary,
		"overdue": overdue,
	})
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

// DebugSliceUsage демонстрирует работу со слайсами
func (h *TaskHandler) DebugSliceUsage(c *gin.Context) {
	tasks := []models.Task{
		{Title: "Задача 1"},
		{Title: "Задача 2"},
	}

	newTask := models.Task{Title: "Задача 3"}
	tasks = append(tasks, newTask)

	c.JSON(http.StatusOK, tasks)
}

// DebugMapUsage демонстрирует работу с мапой
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

// AsyncProcessExample демонстрирует каналы, селекты и контекст
func (h *TaskHandler) AsyncProcessExample(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resultCh := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		resultCh <- "Результат из горутины 1"
	}()

	go func() {
		time.Sleep(3 * time.Second)
		resultCh <- "Результат из горутины 2"
	}()

	select {
	case res := <-resultCh:
		c.JSON(http.StatusOK, gin.H{"result": res})
	case <-ctx.Done():
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "Превышено время ожидания"})
	}
}

func (h *TaskHandler) EnqueueTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	var input struct {
		Title string `json:"title"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := TaskMessage{
		UserID: userID.(uint),
		Title:  input.Title,
	}

	Queue <- task // отправка в очередь

	c.JSON(http.StatusOK, gin.H{"message": "Задача поставлена в очередь"})
}

// ---------- Новый хендлер с sqlc ----------

type TaskSQLCHandler struct {
	Queries *db.Queries
}

func NewTaskSQLCHandler(pool *pgxpool.Pool) *TaskSQLCHandler {
	return &TaskSQLCHandler{
		Queries: db.New(pool),
	}
}

func (h *TaskSQLCHandler) CreateTask(c *gin.Context) {
	var input struct {
		Title     string `json:"title" binding:"required"`
		Completed bool   `json:"completed"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	task, err := h.Queries.CreateTask(ctx, db.CreateTaskParams{
		Title:     input.Title,
		Completed: input.Completed,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания задачи"})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (h *TaskSQLCHandler) GetTasks(c *gin.Context) {
	ctx := context.Background()
	tasks, err := h.Queries.ListTasks(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения задач"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskSQLCHandler) GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}
	ctx := context.Background()
	task, err := h.Queries.GetTask(ctx, int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}
	c.JSON(http.StatusOK, task)
}
