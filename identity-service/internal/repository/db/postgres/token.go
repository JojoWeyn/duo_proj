package postgres

import (
	"context"
	"time"

	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"
	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

func (r *TokenRepository) BlacklistToken(ctx context.Context, token *entity.BlacklistedToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *TokenRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.BlacklistedToken{}).
		Where("token = ? ", token).
		Count(&count).Error
	return count > 0, err
}

func (r *TokenRepository) CleanupExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&entity.BlacklistedToken{}).Error
}
