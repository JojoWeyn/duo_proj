package entity

import (
	"github.com/google/uuid"
	"time"
)

type Progress struct {
	UUID        uuid.UUID `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	UserUUID    uuid.UUID `json:"user_uuid" gorm:"type:uuid"`
	EntityType  string    `json:"entity_type"`
	EntityUUID  uuid.UUID `json:"exercise_uuid" gorm:"type:uuid"`
	Points      int       `json:"points"`
	CompletedAt time.Time `json:"completed_at"`
}

type QuestionAttempt struct {
	UUID         uuid.UUID `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	UserUUID     uuid.UUID `json:"user_uuid" gorm:"type:uuid"`
	QuestionUUID uuid.UUID `json:"question_uuid" gorm:"type:uuid"`
	IsCorrect    bool      `json:"is_correct"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewProgress(userUUID uuid.UUID, entityType string, entityUUID uuid.UUID, points int, completedAt time.Time) *Progress {
	return &Progress{
		UUID:        uuid.New(),
		UserUUID:    userUUID,
		EntityType:  entityType,
		EntityUUID:  entityUUID,
		Points:      points,
		CompletedAt: completedAt,
	}
}

func NewQuestionAttempt(userUUID, questionUUID uuid.UUID, isCorrect bool, createdAt time.Time) *QuestionAttempt {
	return &QuestionAttempt{
		UUID:         uuid.New(),
		UserUUID:     userUUID,
		QuestionUUID: questionUUID,
		IsCorrect:    isCorrect,
		CreatedAt:    createdAt,
	}
}
