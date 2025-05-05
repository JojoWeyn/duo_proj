package dto

import (
	"time"

	"github.com/google/uuid"
)

// ProgressRequestDTO представляет запрос на добавление прогресса
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

// ProgressResponseDTO представляет ответ с прогрессом пользователя
type ProgressResponseDTO struct {
	Exercises []ExerciseProgressDTO `json:"exercises"`
	Lessons   []LessonProgressDTO   `json:"lessons"`
	Courses   []CourseProgressDTO   `json:"courses"`
}

// QuestionAttemptDTO представляет данные о попытке ответа на вопрос
type QuestionAttemptDTO struct {
	UUID         uuid.UUID `json:"uuid"`
	UserUUID     uuid.UUID `json:"user_uuid"`
	QuestionUUID uuid.UUID `json:"question_uuid"`
	IsCorrect    bool      `json:"is_correct"`
	CreatedAt    time.Time `json:"created_at"`
}
