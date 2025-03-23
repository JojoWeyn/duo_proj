package usecase

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"log"
	"time"
)

type CourseRepository interface {
	Create(ctx context.Context, course *entity.Course) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Course, error)
	GetAll(ctx context.Context) ([]*entity.Course, error)
	Update(ctx context.Context, course *entity.Course) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Cacher interface {
	Get(ctx context.Context, key string) ([]*entity.Course, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

type CourseUseCase struct {
	repo  CourseRepository
	cache Cacher
}

func NewCourseUseCase(repo CourseRepository, cache Cacher) *CourseUseCase {
	return &CourseUseCase{
		repo:  repo,
		cache: cache,
	}
}

func (c *CourseUseCase) CreateCourse(ctx context.Context, title, description string, typeID, difficultyID int) error {
	course := entity.NewCourse(title, description, typeID, difficultyID)
	return c.repo.Create(ctx, course)
}

func (c *CourseUseCase) GetCourseByID(ctx context.Context, id uuid.UUID) (*entity.Course, error) {
	return c.repo.GetByID(ctx, id)
}

func (c *CourseUseCase) GetAllCourses(ctx context.Context) ([]*entity.Course, error) {
	cacheKey := "all_courses"

	cachedCourses, err := c.cache.Get(ctx, cacheKey)
	if err == nil && cachedCourses != nil {
		log.Println("Data fetched from cache")
		return cachedCourses, nil
	}

	allCourses, err := c.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	err = c.cache.Set(ctx, cacheKey, allCourses, 10*time.Minute)
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	}

	return allCourses, nil
}

func (c *CourseUseCase) UpdateCourse(ctx context.Context, course *entity.Course) error {
	return c.repo.Update(ctx, course)
}

func (c *CourseUseCase) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	return c.repo.Delete(ctx, id)
}
