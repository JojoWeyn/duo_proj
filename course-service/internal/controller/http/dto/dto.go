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
	UUID         uuid.UUID     `json:"uuid"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	DifficultyID int           `json:"difficulty_id"`
	Difficulty   DifficultyDTO `json:"difficulty"`
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

type StartAttemptResponseDTO struct {
	SessionID uuid.UUID `json:"session_id"`
}

type SubmitAnswerRequestDTO struct {
	QuestionUUID uuid.UUID   `json:"question_id"`
	Answer       interface{} `json:"answer"`
}

type SubmitAnswerResponseDTO struct {
	IsCorrect bool `json:"is_correct"`
}

type FinishAttemptResponseDTO struct {
	IsFinished     bool `json:"is_finished"`
	CorrectAnswers int  `json:"correct_answers"`
	Questions      int  `json:"questions"`
}

type ExerciseAllDTO struct {
	UUID        uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Points      int       `json:"points"`
	Order       int       `json:"order"`
	LessonUUID  uuid.UUID `json:"lesson_uuid"`
}

type ExerciseOneDTO struct {
	UUID          uuid.UUID         `json:"uuid"`
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	Points        int               `json:"points"`
	Order         int               `json:"order"`
	LessonUUID    uuid.UUID         `json:"lesson_uuid"`
	ExerciseFiles []ExerciseFileDTO `json:"exercise_files"`
}

type ExerciseFileDTO struct {
	UUID         uuid.UUID `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	Title        string    `json:"title"`
	FileURL      string    `json:"file_url"`
	ExerciseUUID uuid.UUID `json:"exercise_uuid" gorm:"type:uuid;index"`
}

type LessonsDTO struct {
	UUID         uuid.UUID     `json:"uuid"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	DifficultyID int           `json:"difficulty_id"`
	Difficulty   DifficultyDTO `json:"difficulty"`
	Order        int           `json:"order"`
	CourseUUID   uuid.UUID     `json:"course_uuid"`
}

type DifficultyDTO struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type LessonDTO struct {
	UUID         uuid.UUID       `json:"uuid"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	DifficultyID int             `json:"difficulty_id"`
	Difficulty   DifficultyDTO   `json:"difficulty"`
	Order        int             `json:"order"`
	CourseUUID   uuid.UUID       `json:"course_uuid"`
	LessonFiles  []LessonFileDTO `json:"lesson_files"`
}

type LessonFileDTO struct {
	UUID       uuid.UUID `json:"uuid"`
	Title      string    `json:"title"`
	FileURL    string    `json:"file_url"`
	LessonUUID uuid.UUID `json:"lesson_uuid"`
}

type QuestionsDTO struct {
	UUID         uuid.UUID       `json:"uuid"`
	TypeID       int             `json:"type_id"`
	Type         QuestionTypeDTO `json:"type"`
	Text         string          `json:"text"`
	Order        int             `json:"order"`
	ExerciseUUID uuid.UUID       `json:"exercise_uuid"`
}

type QuestionTypeDTO struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}
