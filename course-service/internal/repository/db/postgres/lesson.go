package postgres

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) *LessonRepository {
	return &LessonRepository{
		db: db,
	}
}

func (l *LessonRepository) Create(ctx context.Context, lesson *entity.Lesson) error {
	return l.db.WithContext(ctx).Create(lesson).Error
}

func (l *LessonRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error) {
	var lesson entity.Lesson

	if err := l.db.WithContext(ctx).
		Where("uuid = ?", id).
		First(&lesson).Error; err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (l *LessonRepository) GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]*entity.Lesson, error) {
	var lessons []*entity.Lesson

	if err := l.db.WithContext(ctx).Where("course_uuid = ?", courseID).Find(&lessons).Error; err != nil {
		return nil, err
	}
	return lessons, nil
}

func (l *LessonRepository) Update(ctx context.Context, lesson *entity.Lesson) error {
	return l.db.WithContext(ctx).Save(lesson).Error
}

func (l LessonRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return l.db.WithContext(ctx).Delete(&entity.Lesson{UUID: id}).Error
}
