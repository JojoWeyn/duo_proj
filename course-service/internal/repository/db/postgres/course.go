package postgres

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{
		db: db,
	}
}

func (c *CourseRepository) Create(ctx context.Context, course *entity.Course) error {
	return c.db.WithContext(ctx).Create(course).Error
}

func (c *CourseRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Course, error) {
	var course entity.Course

	if err := c.db.WithContext(ctx).
		Preload("Lessons").
		Where("uuid = ?", id).First(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (c *CourseRepository) GetAll(ctx context.Context) ([]*entity.Course, error) {
	var courses []*entity.Course

	if err := c.db.WithContext(ctx).
		Preload("Difficulty").
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (c *CourseRepository) Update(ctx context.Context, course *entity.Course) error {
	return c.db.WithContext(ctx).Save(course).Error
}

func (c *CourseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return c.db.WithContext(ctx).Delete(&entity.Course{UUID: id}).Error
}
