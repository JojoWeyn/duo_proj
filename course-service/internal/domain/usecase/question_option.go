package usecase

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/JojoWeyn/duo-proj/course-service/internal/repository/db/postgres"
	"github.com/google/uuid"
)

type QuestionOptionRepository interface {
	Create(ctx context.Context, questionOption *entity.QuestionOption) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.QuestionOption, error)
	GetByQuestionID(ctx context.Context, questionID uuid.UUID) ([]entity.QuestionOption, error)
	Update(ctx context.Context, questionOption *entity.QuestionOption) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type QuestionOptionUseCase struct {
	repo QuestionOptionRepository
}

func NewQuestionOptionUseCase(repo *postgres.QuestionOptionRepository) *QuestionOptionUseCase {
	return &QuestionOptionUseCase{
		repo: repo,
	}
}

func (q *QuestionOptionUseCase) CreateQuestionOption(ctx context.Context, text string, isCorrect bool, questionID uuid.UUID) error {
	option := entity.NewQuestionOption(text, isCorrect, questionID)
	return q.repo.Create(ctx, option)
}

func (q *QuestionOptionUseCase) GetOptionByID(ctx context.Context, id uuid.UUID) (*entity.QuestionOption, error) {
	return q.repo.GetByID(ctx, id)
}

func (q *QuestionOptionUseCase) GetOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]entity.QuestionOption, error) {
	return q.repo.GetByQuestionID(ctx, questionID)
}

func (q *QuestionOptionUseCase) UpdateQuestionOption(ctx context.Context, option *entity.QuestionOption) error {
	return q.repo.Update(ctx, option)
}

func (q *QuestionOptionUseCase) DeleteQuestionOption(ctx context.Context, id uuid.UUID) error {
	return q.repo.Delete(ctx, id)
}
