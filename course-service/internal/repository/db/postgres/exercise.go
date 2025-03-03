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

func (e ExerciseRepository) Create(ctx context.Context, course *entity.Exercise) error {
	//TODO implement me
	panic("implement me")
}

func (e ExerciseRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Exercise, error) {
	//TODO implement me
	panic("implement me")
}

func (e ExerciseRepository) GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]entity.Exercise, error) {
	//TODO implement me
	panic("implement me")
}

func (e ExerciseRepository) Update(ctx context.Context, course *entity.Exercise) error {
	//TODO implement me
	panic("implement me")
}

func (e ExerciseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
