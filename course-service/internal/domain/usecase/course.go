package usecase

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
)

type CourseRepository interface {
	Create(ctx context.Context, course *entity.Course) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Course, error)
	GetAll(ctx context.Context) ([]entity.Course, error)
	Update(ctx context.Context, course *entity.Course) error
	Delete(ctx context.Context, id uuid.UUID) error
}
type CourseUseCase struct {
	repo CourseRepository
}

func NewCourseUseCase(repo CourseRepository) *CourseUseCase {
	return &CourseUseCase{
		repo: repo,
	}
}

func (c *CourseUseCase) CreateCourse(ctx context.Context, title, description string, typeID int) error {
	course := entity.NewCourse(title, description, typeID)
	return c.repo.Create(ctx, course)
}

func (c *CourseUseCase) GetCourseByID(ctx context.Context, id uuid.UUID) (*entity.Course, error) {
	return c.repo.GetByID(ctx, id)
}

func (c *CourseUseCase) GetAllCourses(ctx context.Context) ([]entity.Course, error) {
	return c.repo.GetAll(ctx)
}

func (c *CourseUseCase) UpdateCourse(ctx context.Context, course *entity.Course) error {
	return c.repo.Update(ctx, course)
}

func (c *CourseUseCase) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	return c.repo.Delete(ctx, id)
}
