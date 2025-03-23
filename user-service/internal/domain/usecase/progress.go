package usecase

import (
	"context"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/google/uuid"
	"time"
)

type ProgressRepository interface {
	GetProgress(ctx context.Context, userUUID uuid.UUID) ([]*entity.Progress, error)
	Create(ctx context.Context, progress *entity.Progress) error
	GetByEntityUUID(ctx context.Context, userUUID, entityUUID uuid.UUID) (*entity.Progress, error)
}
type ProgressUseCase struct {
	repo ProgressRepository
}

func NewProgressUseCase(repo ProgressRepository) *ProgressUseCase {
	return &ProgressUseCase{
		repo: repo,
	}
}

func (p *ProgressUseCase) GetProgress(ctx context.Context, userID uuid.UUID) ([]*entity.Progress, error) {
	var progresses []*entity.Progress

	progresses, err := p.repo.GetProgress(ctx, userID)
	if err != nil {
		return nil, err
	}

	return progresses, nil
}

func (p *ProgressUseCase) CheckProgress(ctx context.Context, userID, entityUUID uuid.UUID) bool {
	_, err := p.repo.GetByEntityUUID(ctx, userID, entityUUID)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (p *ProgressUseCase) AddProgress(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, points int, createdAt time.Time) error {
	progress := entity.NewProgress(userID, entityType, entityID, points, createdAt)

	if err := p.repo.Create(ctx, progress); err != nil {
		return err
	}

	return nil
}
