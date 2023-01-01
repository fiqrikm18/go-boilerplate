package dao

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id        int64 `gorm:"primaryKey"`
	Name      string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
