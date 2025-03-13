package entity

import (
	"github.com/google/uuid"
	"time"
)

type Attempt struct {
	UUID         uuid.UUID `json:"uuid" gorm:"type:uuid;primaryKey;"`
	UserUUID     uuid.UUID `json:"user_uuid" gorm:"type:uuid;index"`
	QuestionUUID uuid.UUID `json:"question_uuid" gorm:"type:uuid;index"`
	Answer       string    `json:"answer"`
	IsCorrect    bool      `json:"is_correct"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func NewAttempt(userUUID uuid.UUID, questionUUID uuid.UUID, answer string, isCorrect bool) *Attempt {
	return &Attempt{
		UUID:         uuid.New(),
		UserUUID:     userUUID,
		QuestionUUID: questionUUID,
		Answer:       answer,
		IsCorrect:    isCorrect,
		CreatedAt:    time.Now(),
	}
}
