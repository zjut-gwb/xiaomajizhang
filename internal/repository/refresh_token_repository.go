package repository

import (
	"errors"
	"time"

	"github.com/asikeida/xiaomajizhang/internal/model"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *RefreshTokenRepository) FindByHash(hash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := r.db.Where("token_hash = ?", hash).First(&token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *RefreshTokenRepository) RevokeByHash(hash string, revokedAt time.Time) error {
	return r.db.Model(&model.RefreshToken{}).Where("token_hash = ?", hash).Update("revoked_at", revokedAt).Error
}
