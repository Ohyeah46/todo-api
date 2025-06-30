package models

import (
	"testing"
	"time"
)

func TestTaskMethods(t *testing.T) {
	task := Task{
		Title:       "Test Task",
		Description: "Test Description",
		Completed:   false,
	}
	task.CreatedAt = time.Now().AddDate(0, 0, -8) // 8 дней назад

	// Проверяем метод MarkCompleted
	task.MarkCompleted()
	if !task.Completed {
		t.Error("MarkCompleted не установил Completed в true")
	}

	// Проверяем метод ShortSummary
	expectedSummary := "Test Task: Test Description"
	if summary := task.ShortSummary(); summary != expectedSummary {
		t.Errorf("ShortSummary вернул '%s', ожидалось '%s'", summary, expectedSummary)
	}

	// Проверяем метод IsOverdue (задача должна быть просрочена, так как создана 8 дней назад и не была завершена)
	task.Completed = false
	if overdue := task.IsOverdue(); !overdue {
		t.Error("IsOverdue вернул false, ожидалось true")
	}

	// Теперь отмечаем как выполненную — просрочка должна быть false
	task.MarkCompleted()
	if overdue := task.IsOverdue(); overdue {
		t.Error("IsOverdue вернул true для выполненной задачи, ожидалось false")
	}
}
