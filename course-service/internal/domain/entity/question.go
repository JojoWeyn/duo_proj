package entity

import "github.com/google/uuid"

type Question struct {
	UUID            string           `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	TypeID          int              `json:"type"`
	Text            string           `json:"text"`
	Order           int              `json:"order"`
	ExerciseUUID    uuid.UUID        `json:"exercise_uuid" gorm:"type:uuid;index"`
	Exercise        Exercise         `json:"exercise" gorm:"foreignKey:ExerciseUUID;references:UUID"`
	QuestionOptions []QuestionOption `json:"question_options" gorm:"foreignKey:QuestionUUID;constraint:OnDelete:CASCADE"`
}
