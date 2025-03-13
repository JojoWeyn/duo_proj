package postgres

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttemptRepository struct {
	db *gorm.DB
}

func NewAttemptRepository(db *gorm.DB) *AttemptRepository {
	return &AttemptRepository{
		db: db,
	}
}

func (a *AttemptRepository) Create(ctx context.Context, attempt *entity.Attempt) error {
	return a.db.WithContext(ctx).Create(attempt).Error
}

func (a *AttemptRepository) GetAttemptsByUser(ctx context.Context, userUUID uuid.UUID) ([]entity.Attempt, error) {
	var attempts []entity.Attempt
	if err := a.db.WithContext(ctx).Where("user_uuid = ?", userUUID).Find(&attempts).Error; err != nil {
		return nil, err
	}
	return attempts, nil
}
