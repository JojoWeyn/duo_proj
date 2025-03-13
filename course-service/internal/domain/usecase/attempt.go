package usecase

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
)

type AttemptRepository interface {
	Create(ctx context.Context, attempt *entity.Attempt) error
	GetAttemptsByUser(ctx context.Context, userUUID uuid.UUID) ([]entity.Attempt, error)
}

type AttemptUseCase struct {
	repo AttemptRepository
}

func NewAttemptUseCase(repo AttemptRepository) *AttemptUseCase {
	return &AttemptUseCase{
		repo: repo,
	}
}

func (a *AttemptUseCase) CreateAttempt(ctx context.Context, userUUID, questionUUID uuid.UUID, answer string, isCorrect bool) error {
	attempt := entity.NewAttempt(
		userUUID,
		questionUUID,
		answer,
		isCorrect,
	)
	return a.repo.Create(ctx, attempt)
}

func (a *AttemptUseCase) GetAttemptsByUser(ctx context.Context, userUUID uuid.UUID) ([]entity.Attempt, error) {
	return a.repo.GetAttemptsByUser(ctx, userUUID)
}
