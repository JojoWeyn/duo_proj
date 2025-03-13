package usecase

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
)

type ExerciseRepository interface {
	Create(ctx context.Context, course *entity.Exercise) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Exercise, error)
	GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*entity.Exercise, error)
	Update(ctx context.Context, course *entity.Exercise) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ExerciseUsecase struct {
	repo ExerciseRepository
}

func NewExerciseUseCase(repo ExerciseRepository) *ExerciseUsecase {
	return &ExerciseUsecase{
		repo: repo,
	}
}

func (e *ExerciseUsecase) CreateExercise(ctx context.Context, title, description string, points, order int, lessonUUID uuid.UUID) error {
	exercise := entity.NewExercise(title, description, points, order, lessonUUID)
	return e.repo.Create(ctx, exercise)
}

func (e *ExerciseUsecase) GetExerciseByID(ctx context.Context, id uuid.UUID) (*entity.Exercise, error) {
	return e.repo.GetByID(ctx, id)
}

func (e *ExerciseUsecase) GetExercisesByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*entity.Exercise, error) {
	return e.repo.GetByLessonID(ctx, lessonID)
}

func (e *ExerciseUsecase) UpdateExercise(ctx context.Context, exercise *entity.Exercise) error {
	return e.repo.Update(ctx, exercise)
}

func (e *ExerciseUsecase) DeleteExercise(ctx context.Context, id uuid.UUID) error {
	return e.repo.Delete(ctx, id)
}
