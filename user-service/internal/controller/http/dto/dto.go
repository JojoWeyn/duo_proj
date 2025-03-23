package dto

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type UserAchievementsDTO struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Secret       bool      `json:"secret"`
	CurrentCount int       `json:"current_count"`
	Achieved     bool      `json:"achieved"`
	CreatedAt    time.Time `json:"created_at"`
	AchievedAt   time.Time `json:"achieved_at"`
}

type AchievementsDTO struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Condition   json.RawMessage `json:"condition"`
	Secret      bool            `json:"secret"`
	CreatedAt   time.Time       `json:"created_at"`
}

type ExerciseProgressDTO struct {
	UUID         uuid.UUID `json:"uuid" gorm:"type:uuid;primaryKey;"`
	ExerciseUUID uuid.UUID `json:"exercise_uuid" gorm:"type:uuid;index"`
	TotalPoints  int       `json:"total_points"`
	CompletedAt  time.Time `json:"completed_at" gorm:"autoCreateTime"`
}

type LessonProgressDTO struct {
	UUID        uuid.UUID `json:"uuid" gorm:"type:uuid;primaryKey;"`
	LessonUUID  uuid.UUID `json:"lesson_uuid" gorm:"type:uuid;index"`
	TotalPoints int       `json:"total_points"`
	CompletedAt time.Time `json:"completed_at" gorm:"autoCreateTime"`
}

type CourseProgressDTO struct {
	UUID        uuid.UUID `json:"uuid" gorm:"type:uuid;primaryKey;"`
	CourseUUID  uuid.UUID `json:"course_uuid" gorm:"type:uuid;index"`
	TotalPoints int       `json:"total_points"`
	CompletedAt time.Time `json:"completed_at" gorm:"autoCreateTime"`
}
