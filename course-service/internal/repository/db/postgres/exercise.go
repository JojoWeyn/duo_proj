package postgres

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExerciseRepository struct {
	db *gorm.DB
}

func NewExerciseRepository(db *gorm.DB) *ExerciseRepository {
	return &ExerciseRepository{
		db: db,
	}
}

func (e *ExerciseRepository) GetByQuestionID(ctx context.Context, questionID uuid.UUID) (*entity.Exercise, error) {
	var exercise *entity.Exercise
	e.db.
		Joins("JOIN questions ON questions.exercise_uuid = exercises.uuid").
		Where("questions.uuid = ?", questionID).
		First(&exercise)
	return exercise, nil
}

func (e *ExerciseRepository) Create(ctx context.Context, course *entity.Exercise) error {
	return e.db.WithContext(ctx).Create(course).Error
}

func (e *ExerciseRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Exercise, error) {
	var exercise entity.Exercise

	if err := e.db.WithContext(ctx).
		Where("uuid = ?", id).
		Preload("ExerciseFiles").
		First(&exercise).Error; err != nil {
		return nil, err
	}
	return &exercise, nil
}

func (e *ExerciseRepository) GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*entity.Exercise, error) {
	var exercises []*entity.Exercise

	if err := e.db.WithContext(ctx).Where("lesson_uuid = ?", lessonID).Find(&exercises).Error; err != nil {
		return nil, err
	}
	return exercises, nil
}

func (e *ExerciseRepository) Update(ctx context.Context, course *entity.Exercise) error {
	return e.db.WithContext(ctx).Save(course).Error
}

func (e *ExerciseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return e.db.WithContext(ctx).Delete(&entity.Exercise{UUID: id}).Error
}

func (e *ExerciseRepository) AddFile(ctx context.Context, file *entity.ExerciseFile) error {
	return e.db.WithContext(ctx).Create(file).Error
}
