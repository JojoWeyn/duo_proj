package usecase

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
)

type MatchingPairRepository interface {
	Create(ctx context.Context, questionMatchingPair *entity.MatchingPair) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.MatchingPair, error)
	GetByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*entity.MatchingPair, error)
	Update(ctx context.Context, questionMatchingPair *entity.MatchingPair) error
}
type MatchingPairUseCase struct {
	repo MatchingPairRepository
}

func NewMatchingPairUseCase(repo MatchingPairRepository) *MatchingPairUseCase {
	return &MatchingPairUseCase{
		repo: repo,
	}
}

func (m *MatchingPairUseCase) CreateMatchingPair(ctx context.Context, left, right string, questionUUID uuid.UUID) error {
	questionMatchingPair := entity.NewMatchingPair(left, right, questionUUID)
	return m.repo.Create(ctx, questionMatchingPair)
}

func (m *MatchingPairUseCase) DeleteMatchingPair(ctx context.Context, id uuid.UUID) error {
	return m.repo.Delete(ctx, id)
}

func (m *MatchingPairUseCase) GetMatchingPairByID(ctx context.Context, id uuid.UUID) (*entity.MatchingPair, error) {
	return m.repo.GetByID(ctx, id)
}

func (m *MatchingPairUseCase) GetMatchingPairsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*entity.MatchingPair, error) {
	return m.repo.GetByQuestionID(ctx, questionID)
}

func (m *MatchingPairUseCase) UpdateMatchingPair(ctx context.Context, questionMatchingPair *entity.MatchingPair) error {
	return m.repo.Update(ctx, questionMatchingPair)
}
