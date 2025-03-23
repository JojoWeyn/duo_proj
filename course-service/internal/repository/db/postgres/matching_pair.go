package postgres

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MatchingPairRepository struct {
	db *gorm.DB
}

func NewMatchingPairRepository(db *gorm.DB) *MatchingPairRepository {
	return &MatchingPairRepository{
		db: db,
	}
}

func (q *MatchingPairRepository) Create(ctx context.Context, questionMatchingPair *entity.MatchingPair) error {
	return q.db.WithContext(ctx).Create(questionMatchingPair).Error
}

func (q *MatchingPairRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return q.db.WithContext(ctx).Where("uuid = ?", id).Delete(&entity.MatchingPair{}).Error
}

func (q *MatchingPairRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.MatchingPair, error) {
	var questionMatchingPair entity.MatchingPair
	if err := q.db.WithContext(ctx).Where("uuid = ?", id).First(&questionMatchingPair).Error; err != nil {
		return nil, err
	}
	return &questionMatchingPair, nil
}

func (q *MatchingPairRepository) GetByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*entity.MatchingPair, error) {
	var questionMatchingPairs []*entity.MatchingPair
	if err := q.db.WithContext(ctx).Where("question_uuid = ?", questionID).Find(&questionMatchingPairs).Error; err != nil {
		return nil, err
	}
	return questionMatchingPairs, nil
}

func (q *MatchingPairRepository) Update(ctx context.Context, questionMatchingPair *entity.MatchingPair) error {
	return q.db.WithContext(ctx).Save(questionMatchingPair).Error
}
