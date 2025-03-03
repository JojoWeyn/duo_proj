package entity

import "github.com/google/uuid"

type QuestionOption struct {
	UUID         uuid.UUID `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	Text         string    `json:"text"`
	IsCorrect    bool      `json:"is_correct"`
	QuestionUUID uuid.UUID `json:"question_uuid" gorm:"type:uuid;index"`
	Question     Question  `json:"question" gorm:"foreignKey:QuestionUUID;references:UUID"`
}

func NewQuestionOption(text string, isCorrect bool, questionUUID uuid.UUID) *QuestionOption {
	return &QuestionOption{
		UUID:         uuid.New(),
		Text:         text,
		IsCorrect:    isCorrect,
		QuestionUUID: questionUUID,
	}
}
