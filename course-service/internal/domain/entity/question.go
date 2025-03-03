package entity

import "github.com/google/uuid"

type Question struct {
	UUID            uuid.UUID        `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	TypeID          int              `json:"type_id" gorm:"index"`
	QuestionType    QuestionType     `json:"type" gorm:"foreignKey:TypeID;references:ID"`
	Text            string           `json:"text"`
	Order           int              `json:"order"`
	ExerciseUUID    uuid.UUID        `json:"exercise_uuid" gorm:"type:uuid;index"`
	Exercise        Exercise         `json:"exercise" gorm:"foreignKey:ExerciseUUID;references:UUID"`
	QuestionOptions []QuestionOption `json:"question_options" gorm:"foreignKey:QuestionUUID;constraint:OnDelete:CASCADE"`
}

func NewQuestion(text string, typeID, order int, exerciseUUID uuid.UUID) *Question {
	return &Question{
		UUID:            uuid.New(),
		TypeID:          typeID,
		Text:            text,
		Order:           order,
		ExerciseUUID:    exerciseUUID,
		QuestionOptions: []QuestionOption{},
	}
}
