package database

import (
	"github.com/asikeida/xiaomajizhang/internal/model"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
	)
}
