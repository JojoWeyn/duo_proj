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

func (l LessonRepository) Create(ctx context.Context, lesson *entity.Lesson) error {
	//TODO implement me
	panic("implement me")
}

func (l LessonRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error) {
	//TODO implement me
	panic("implement me")
}

func (l LessonRepository) GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]entity.Lesson, error) {
	//TODO implement me
	panic("implement me")
}

func (l LessonRepository) Update(ctx context.Context, lesson *entity.Lesson) error {
	//TODO implement me
	panic("implement me")
}

func (l LessonRepository) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
