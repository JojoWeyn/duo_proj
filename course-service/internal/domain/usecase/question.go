package usecase

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
)

type QuestionRepository interface {
	Create(ctx context.Context, question *entity.Question) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Question, error)
	GetByExerciseID(ctx context.Context, exerciseID uuid.UUID) ([]entity.Question, error)
	Update(ctx context.Context, question *entity.Question) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type QuestionUseCase struct {
	repo QuestionRepository
}

func NewQuestionUseCase(repo QuestionRepository) *QuestionUseCase {
	return &QuestionUseCase{
		repo: repo,
	}
}

func (q *QuestionUseCase) CreateQuestion(ctx context.Context, text string, typeID, order int, exerciseUUID uuid.UUID) error {
	question := entity.NewQuestion(text, typeID, order, exerciseUUID)
	return q.repo.Create(ctx, question)
}

func (q *QuestionUseCase) GetQuestionByID(ctx context.Context, id uuid.UUID) (*entity.Question, error) {
	return q.repo.GetByID(ctx, id)
}

func (q *QuestionUseCase) GetQuestionsByExerciseID(ctx context.Context, exerciseID uuid.UUID) ([]entity.Question, error) {
	return q.repo.GetByExerciseID(ctx, exerciseID)
}

func (q *QuestionUseCase) UpdateQuestion(ctx context.Context, question *entity.Question) error {
	return q.repo.Update(ctx, question)
}

func (q *QuestionUseCase) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	return q.repo.Delete(ctx, id)
}
