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

func (q *QuestionOptionRepository) Create(ctx context.Context, questionOption *entity.QuestionOption) error {
	return q.db.WithContext(ctx).Create(questionOption).Error
}

func (q *QuestionOptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.QuestionOption, error) {
	var questionOption entity.QuestionOption
	if err := q.db.WithContext(ctx).Where("uuid = ?", id).First(&questionOption).Error; err != nil {
		return nil, err
	}
	return &questionOption, nil
}

func (q *QuestionOptionRepository) GetByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*entity.QuestionOption, error) {
	var questionOptions []*entity.QuestionOption
	if err := q.db.WithContext(ctx).Where("question_uuid = ?", questionID).Find(&questionOptions).Error; err != nil {
		return nil, err
	}
	return questionOptions, nil
}

func (q *QuestionOptionRepository) Update(ctx context.Context, questionOption *entity.QuestionOption) error {
	return q.db.WithContext(ctx).Save(questionOption).Error
}

func (q *QuestionOptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return q.db.WithContext(ctx).Where("uuid = ?", id).Delete(&entity.QuestionOption{}).Error
}
