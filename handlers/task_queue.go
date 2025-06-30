package handlers

import (
	"log"
	"time"
)

type TaskMessage struct {
	UserID uint
	Title  string
}

var Queue chan TaskMessage

func InitTaskWorker() {
	Queue = make(chan TaskMessage, 100)

	go func() {
		for task := range Queue {
			log.Printf("[WORKER] Обработка задачи: user=%d title=%s\n", task.UserID, task.Title)
			time.Sleep(2 * time.Second)
			log.Printf("[WORKER] ✅ Готово: %s\n", task.Title)
		}
	}()
}
