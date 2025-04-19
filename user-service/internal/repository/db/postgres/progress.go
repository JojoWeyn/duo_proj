package postgres

import (
	"context"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProgressRepository struct {
	db *gorm.DB
}

func NewProgressRepository(db *gorm.DB) *ProgressRepository {
	return &ProgressRepository{
		db: db,
	}
}

func (p *ProgressRepository) Create(ctx context.Context, progress *entity.Progress) error {
	return p.db.WithContext(ctx).Create(progress).Error
}

func (p *ProgressRepository) GetByEntityUUID(ctx context.Context, userUUID, entityUUID uuid.UUID) (*entity.Progress, error) {
	var progress entity.Progress
	err := p.db.WithContext(ctx).
		Where("user_uuid = ? AND entity_uuid = ?", userUUID, entityUUID).
		First(&progress).Error
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

func (p *ProgressRepository) GetProgress(ctx context.Context, userUUID uuid.UUID) ([]*entity.Progress, error) {
	var progresses []*entity.Progress

	if err := p.db.WithContext(ctx).Where("user_uuid = ?", userUUID).Find(&progresses).Error; err != nil {
		return nil, err
	}

	return progresses, nil
}

func (p *ProgressRepository) Update(ctx context.Context, progress *entity.Progress) error {
	return p.db.WithContext(ctx).
		Model(&entity.Progress{}).
		Where("user_uuid = ? AND entity_uuid = ?", progress.UserUUID, progress.EntityUUID).
		Updates(progress).Error
}

func (p *ProgressRepository) Delete(ctx context.Context, userUUID, entityUUID uuid.UUID) error {
	return p.db.WithContext(ctx).
		Where("user_uuid = ? AND entity_uuid = ?", userUUID, entityUUID).
		Delete(&entity.Progress{}).Error
}
