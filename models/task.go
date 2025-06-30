package models

import "gorm.io/gorm"

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
