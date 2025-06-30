package models

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	UserID      uint   `json:"user_id"` // Ссылка на пользователя
}

func (t *Task) SetTitle(newTitle string) {
	t.Title = newTitle
}

// Метод отмечает задачу как выполненную
func (t *Task) MarkCompleted() {
	t.Completed = true
}

// Метод возвращает короткое описание задачи
func (t *Task) ShortSummary() string {
	return t.Title + ": " + t.Description
}

// Метод проверяет, просрочена ли задача (например, если CreatedAt старше 7 дней)
func (t *Task) IsOverdue() bool {
	return time.Since(t.CreatedAt) > 7*24*time.Hour && !t.Completed
}
