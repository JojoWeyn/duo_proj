package usecase

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
)

type LessonRepository interface {
	Create(ctx context.Context, lesson *entity.Lesson) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error)
	GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]entity.Lesson, error)
	Update(ctx context.Context, lesson *entity.Lesson) error
	Delete(ctx context.Context, id uuid.UUID) error
}
type LessonUseCase struct {
	repo LessonRepository
}

func NewLessonUseCase(repo LessonRepository) *LessonUseCase {
	return &LessonUseCase{
		repo: repo,
	}
}

func (l *LessonUseCase) CreateLesson(ctx context.Context, title, description string, difficultyID, order int, courseUUID uuid.UUID) error {
	lesson := entity.NewLesson(title, description, difficultyID, order, courseUUID)
	return l.repo.Create(ctx, lesson)
}

func (l *LessonUseCase) GetLessonByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error) {
	return l.repo.GetByID(ctx, id)
}

func (l *LessonUseCase) GetLessonsByCourseID(ctx context.Context, courseID uuid.UUID) ([]entity.Lesson, error) {
	return l.repo.GetByCourseID(ctx, courseID)
}

func (l *LessonUseCase) UpdateLesson(ctx context.Context, lesson *entity.Lesson) error {
	return l.repo.Update(ctx, lesson)
}

func (l *LessonUseCase) DeleteLesson(ctx context.Context, id uuid.UUID) error {
	return l.repo.Delete(ctx, id)
}
