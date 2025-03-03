package postgres

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionOptionRepository struct {
	db *gorm.DB
}

func NewQuestionOptionRepository(db *gorm.DB) *QuestionOptionRepository {
	return &QuestionOptionRepository{
		db: db,
	}
}

func (q QuestionOptionRepository) Create(ctx context.Context, questionOption *entity.QuestionOption) error {
	//TODO implement me
	panic("implement me")
}

func (q QuestionOptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.QuestionOption, error) {
	//TODO implement me
	panic("implement me")
}

func (q QuestionOptionRepository) GetByQuestionID(ctx context.Context, questionID uuid.UUID) ([]entity.QuestionOption, error) {
	//TODO implement me
	panic("implement me")
}

func (q QuestionOptionRepository) Update(ctx context.Context, questionOption *entity.QuestionOption) error {
	//TODO implement me
	panic("implement me")
}

func (q QuestionOptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
