package postgres

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Completion struct {
	UserUUID   uuid.UUID `gorm:"primaryKey"`
	EntityUUID uuid.UUID `gorm:"primaryKey"`
	EntityType string    `gorm:"index"` // "exercise", "lesson", "course"
	Points     int
}

type CompletionRepo struct {
	db *gorm.DB
}

func NewCompletionRepo(db *gorm.DB) *CompletionRepo {
	return &CompletionRepo{db: db}
}

func (c *CompletionRepo) MarkExerciseCompleted(ctx context.Context, userUUID, exerciseUUID uuid.UUID, points int) error {
	return c.markCompleted(ctx, userUUID, exerciseUUID, "exercise", points)
}

func (c *CompletionRepo) IsExerciseCompleted(ctx context.Context, userUUID, exerciseUUID uuid.UUID) (bool, int, error) {
	return c.isCompleted(ctx, userUUID, exerciseUUID, "exercise")
}

func (c *CompletionRepo) MarkLessonCompleted(ctx context.Context, userUUID, lessonUUID uuid.UUID, points int) error {
	return c.markCompleted(ctx, userUUID, lessonUUID, "lesson", points)
}

func (c *CompletionRepo) IsLessonCompleted(ctx context.Context, userUUID, lessonUUID uuid.UUID) (bool, int, error) {
	return c.isCompleted(ctx, userUUID, lessonUUID, "lesson")
}

func (c *CompletionRepo) MarkCourseCompleted(ctx context.Context, userUUID, courseUUID uuid.UUID, points int) error {
	return c.markCompleted(ctx, userUUID, courseUUID, "course", points)
}

func (c *CompletionRepo) IsCourseCompleted(ctx context.Context, userUUID, courseUUID uuid.UUID) (bool, int, error) {
	return c.isCompleted(ctx, userUUID, courseUUID, "course")
}

func (c *CompletionRepo) markCompleted(ctx context.Context, userUUID, entityUUID uuid.UUID, entityType string, points int) error {
	completion := Completion{
		UserUUID:   userUUID,
		EntityUUID: entityUUID,
		EntityType: entityType,
		Points:     points,
	}
	return c.db.WithContext(ctx).Create(&completion).Error
}

func (c *CompletionRepo) isCompleted(ctx context.Context, userUUID, entityUUID uuid.UUID, entityType string) (bool, int, error) {
	var completion Completion
	result := c.db.WithContext(ctx).Where("user_uuid = ? AND entity_uuid = ? AND entity_type = ?", userUUID, entityUUID, entityType).First(&completion)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, 0, nil
		}
		return false, 0, result.Error
	}
	return true, completion.Points, nil
}
