package database

import (
	"fmt"

	"github.com/asikeida/xiaomajizhang/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQL(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect mysql: %w", err)
	}
	return db, nil
}
