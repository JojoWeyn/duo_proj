package postgres

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionMatchingPairRepository struct {
	db *gorm.DB
}

func NewQuestionMatchingPairRepository(db *gorm.DB) *QuestionMatchingPairRepository {
	return &QuestionMatchingPairRepository{
		db: db,
	}
}

func (q *QuestionMatchingPairRepository) Create(ctx context.Context, questionMatchingPair *entity.MatchingPair) error {
	return q.db.WithContext(ctx).Create(questionMatchingPair).Error
}

func (q *QuestionMatchingPairRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return q.db.WithContext(ctx).Where("uuid = ?", id).Delete(&entity.MatchingPair{}).Error
}

func (q *QuestionMatchingPairRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.MatchingPair, error) {
	var questionMatchingPair entity.MatchingPair
	if err := q.db.WithContext(ctx).Where("uuid = ?", id).First(&questionMatchingPair).Error; err != nil {
		return nil, err
	}
	return &questionMatchingPair, nil
}

func (q *QuestionMatchingPairRepository) GetByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*entity.MatchingPair, error) {
	var questionMatchingPairs []*entity.MatchingPair
	if err := q.db.WithContext(ctx).Where("question_uuid = ?", questionID).Find(&questionMatchingPairs).Error; err != nil {
		return nil, err
	}
	return questionMatchingPairs, nil
}

func (q *QuestionMatchingPairRepository) Update(ctx context.Context, questionMatchingPair *entity.MatchingPair) error {
	return q.db.WithContext(ctx).Save(questionMatchingPair).Error
}
