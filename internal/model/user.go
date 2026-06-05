package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	Username     string `gorm:"size:50;uniqueIndex;not null"`
	Email        string `gorm:"size:100;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	Nickname     string `gorm:"size:50;not null"`
	AvatarURL    string `gorm:"size:500"`
	Status       int8   `gorm:"not null;default:1"`
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
