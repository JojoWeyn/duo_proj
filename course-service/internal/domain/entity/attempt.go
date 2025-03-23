package entity

import (
	"github.com/google/uuid"
	"time"
)

type Attempt struct {
	UUID               uuid.UUID `json:"uuid" gorm:"type:uuid;primaryKey;"`
	UserUUID           uuid.UUID `json:"user_uuid" gorm:"type:uuid;index"`
	QuestionUUID       uuid.UUID `json:"question_uuid" gorm:"type:uuid;index"`
	AttemptSessionUUID uuid.UUID `json:"attempt_session_uuid" gorm:"type:uuid;index"`
	Answer             string    `json:"answer"`
	IsCorrect          bool      `json:"is_correct"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type AttemptSession struct {
	UUID         uuid.UUID  `json:"uuid" gorm:"type:uuid;primaryKey;"`
	UserUUID     uuid.UUID  `json:"user_uuid" gorm:"type:uuid;index"`
	ExerciseUUID uuid.UUID  `json:"exercise_uuid" gorm:"type:uuid;index"`
	StartedAt    time.Time  `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at"`
}

func NewAttemptSession(userUUID, exerciseUUID uuid.UUID) *AttemptSession {
	return &AttemptSession{
		UUID:         uuid.New(),
		UserUUID:     userUUID,
		ExerciseUUID: exerciseUUID,
		StartedAt:    time.Now(),
		FinishedAt:   nil,
	}
}

func NewAttempt(userUUID, sessionUUID, questionUUID uuid.UUID, answer string, isCorrect bool) *Attempt {
	return &Attempt{
		UUID:               uuid.New(),
		UserUUID:           userUUID,
		QuestionUUID:       questionUUID,
		AttemptSessionUUID: sessionUUID,
		Answer:             answer,
		IsCorrect:          isCorrect,
		CreatedAt:          time.Now(),
	}
}
