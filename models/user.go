package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password"` // Будет храниться хеш пароля
	Tasks    []Task
}
