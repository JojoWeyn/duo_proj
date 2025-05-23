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

func (q *QuestionRepository) Create(ctx context.Context, question *entity.Question) error {
	return q.db.WithContext(ctx).Create(question).Error
}

func (q *QuestionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Question, error) {
	var question entity.Question
	if err := q.db.WithContext(ctx).
		Preload("QuestionType").
		Preload("QuestionImages").
		Preload("MatchingPairs").
		Preload("QuestionOptions").Where("uuid = ?", id).First(&question).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

func (q *QuestionRepository) GetByExerciseID(ctx context.Context, exerciseID uuid.UUID) ([]entity.Question, error) {
	var questions []entity.Question
	if err := q.db.WithContext(ctx).
		Preload("QuestionType").
		Where("exercise_uuid = ?", exerciseID).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (q *QuestionRepository) Update(ctx context.Context, question *entity.Question) error {
	// Только изменённые поля
	return q.db.WithContext(ctx).
		Model(&entity.Question{}).
		Where("uuid = ?", question.UUID).
		Updates(map[string]interface{}{
			"text":    question.Text,
			"type_id": question.TypeID,
			"order":   question.Order,
		}).Error
}

func (q *QuestionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return q.db.WithContext(ctx).Where("uuid = ?", id).Delete(&entity.Question{}).Error
}

func (q *QuestionRepository) GetCountByExerciseID(ctx context.Context, exerciseID uuid.UUID) (int64, error) {
	var totalQuestions int64
	if err := q.db.Model(&entity.Question{}).Where("exercise_uuid = ?", exerciseID).Count(&totalQuestions).Error; err != nil {
		return 0, err
	}
	return totalQuestions, nil
}

func (q *QuestionRepository) AddImage(ctx context.Context, file *entity.QuestionImage) error {
	return q.db.WithContext(ctx).Create(file).Error
}

func (q *QuestionRepository) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return q.db.WithContext(ctx).Delete(&entity.QuestionImage{UUID: id}).Error
}
