package dto

import (
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
)

type CourseCreateRequestDTO struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	TypeID       int    `json:"type_id"`
	DifficultyID int    `json:"difficulty_id"`
}

type LessonCreateRequestDTO struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	DifficultyID int       `json:"difficulty_id"`
	Order        int       `json:"order"`
	CourseUUID   uuid.UUID `json:"courseUUID"`
}

type ExerciseCreateRequestDTO struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Points      int       `json:"points"`
	Order       int       `json:"order"`
	LessonUUID  uuid.UUID `json:"lessonUUID"`
}

type QuestionCreateRequestDTO struct {
	TypeID       int       `json:"type_id"`
	Text         string    `json:"text"`
	Order        int       `json:"order"`
	ExerciseUUID uuid.UUID `json:"exercise_uuid"`
}

type MatchingPairCreateRequestDTO struct {
	LeftText     string `json:"left_text"`
	RightText    string `json:"right_text"`
	QuestionUUID uuid.UUID
}

type QuestionOptionCreateRequestDTO struct {
	Text         string    `json:"text"`
	IsCorrect    bool      `json:"is_correct"`
	QuestionUUID uuid.UUID `json:"questionUUID"`
}
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
	Images          []QuestionImageDTO  `json:"images"`
	Order           int                 `json:"order"`
	ExerciseUUID    uuid.UUID           `json:"exercise_uuid"`
	QuestionOptions []QuestionOptionDTO `json:"question_options"`
	Matching        QuestionMatchingDTO `json:"matching"`
}

type QuestionImageDTO struct {
	ImageURL string `json:"image_url"`
}

type QuestionOptionDTO struct {
	UUID uuid.UUID `json:"uuid"`
	Text string    `json:"text"`
}

type QuestionMatchingDTO struct {
	LeftSide  []string `json:"left_side"`
	RightSide []string `json:"right_side"`
}
