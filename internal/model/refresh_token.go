package model

import "time"

type RefreshToken struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	UserID    uint64    `gorm:"index;not null"`
	TokenHash string    `gorm:"size:255;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"index;not null"`
	RevokedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
