package entity

import "github.com/google/uuid"

type MatchingPair struct {
	UUID         uuid.UUID `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	LeftText     string    `json:"left_text"`
	RightText    string    `json:"right_text"`
	QuestionUUID uuid.UUID `json:"question_uuid" gorm:"type:uuid;index"`
}

func NewMatchingPair(left, right string, questionUUID uuid.UUID) *MatchingPair {
	return &MatchingPair{
		UUID:         uuid.New(),
		LeftText:     left,
		RightText:    right,
		QuestionUUID: questionUUID,
	}
}
