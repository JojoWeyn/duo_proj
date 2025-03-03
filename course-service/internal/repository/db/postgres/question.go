package postgres

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{
		db: db,
	}
}

func (q QuestionRepository) Create(ctx context.Context, question *entity.Question) error {
	//TODO implement me
	panic("implement me")
}

func (q QuestionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Question, error) {
	//TODO implement me
	panic("implement me")
}

func (q QuestionRepository) GetByExerciseID(ctx context.Context, exerciseID uuid.UUID) ([]entity.Question, error) {
	//TODO implement me
	panic("implement me")
}

func (q QuestionRepository) Update(ctx context.Context, question *entity.Question) error {
	//TODO implement me
	panic("implement me")
}

func (q QuestionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
