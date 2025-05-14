package usecase

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"log"
	"strings"
	"time"
)

type CourseRepository interface {
	Create(ctx context.Context, course *entity.Course) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Course, error)
	GetAll(ctx context.Context) ([]entity.Course, error)
	Update(ctx context.Context, course *entity.Course) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

type CourseUseCase struct {
	repo  CourseRepository
	cache Cache
}

func NewCourseUseCase(repo CourseRepository, cache Cache) *CourseUseCase {
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

func (c *CourseUseCase) GetAllCourses(ctx context.Context, title string, diffId int, typeId int) ([]entity.Course, error) {
	cacheKey := "all_courses"

	var allCourses []entity.Course
	if err := c.cache.Get(ctx, cacheKey, &allCourses); err == nil && allCourses != nil {
		log.Println("Data fetched from cache")
		return filterCourses(allCourses, title, diffId, typeId), nil
	}

	allCourses, err := c.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	err = c.cache.Set(ctx, cacheKey, allCourses, 1*time.Minute)
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	}

	return filterCourses(allCourses, title, diffId, typeId), nil
}

func filterCourses(courses []entity.Course, title string, diffId, typeId int) []entity.Course {
	filtered := make([]entity.Course, 0)

	for _, course := range courses {
		if typeId == 2 && course.TypeID != 2 {
			continue
		}
		if title != "" && !strings.HasPrefix(strings.ToLower(course.Title), strings.ToLower(title)) {
			continue
		}
		if diffId != 0 && course.DifficultyID != diffId {
			continue
		}
		filtered = append(filtered, course)
	}

	return filtered
}

func (c *CourseUseCase) UpdateCourse(ctx context.Context, course *entity.Course) error {
	return c.repo.Update(ctx, course)
}

func (c *CourseUseCase) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	return c.repo.Delete(ctx, id)
}
