package dto

import (
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
)

type CourseSmallDTO struct {
	UUID         uuid.UUID         `json:"uuid"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	DifficultyID int               `json:"difficulty_id"`
	Difficulty   entity.Difficulty `json:"difficulty"`
}

type CourseInfoDTO struct {
	UUID         uuid.UUID         `json:"uuid"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	TypeID       int               `json:"type_id"`
	CourseType   entity.CourseType `json:"course_type"`
	DifficultyID int               `json:"difficulty_id"`
	Difficulty   entity.Difficulty `json:"difficulty"`
}

type QuestionDTO struct {
	UUID            uuid.UUID           `json:"uuid"`
	TypeID          int                 `json:"type_id"`
	QuestionType    entity.QuestionType `json:"type"`
	Text            string              `json:"text"`
	Order           int                 `json:"order"`
	ExerciseUUID    uuid.UUID           `json:"exercise_uuid"`
	QuestionOptions []QuestionOptionDTO `json:"question_options"`
	Matching        QuestionMatchingDTO `json:"matching"`
}

type QuestionOptionDTO struct {
	UUID uuid.UUID `json:"uuid"`
	Text string    `json:"text"`
}

type QuestionMatchingDTO struct {
	LeftSide  []string `json:"left_side"`
	RightSide []string `json:"right_side"`
}
