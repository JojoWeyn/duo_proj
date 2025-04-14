package entity

import "github.com/google/uuid"

type Lesson struct {
	UUID         uuid.UUID    `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	Title        string       `json:"title"`
	Description  string       `json:"description" gorm:"type:text"`
	DifficultyID int          `json:"difficulty_id" gorm:"index"`
	Difficulty   Difficulty   `json:"difficulty" gorm:"foreignKey:DifficultyID;references:ID"`
	Order        int          `json:"order"`
	CourseUUID   uuid.UUID    `json:"course_uuid" gorm:"type:uuid;index"`
	Exercises    []Exercise   `json:"exercises" gorm:"foreignKey:LessonUUID;constraint:OnDelete:CASCADE"`
	LessonFiles  []LessonFile `json:"lesson_files" gorm:"foreignKey:LessonUUID;constraint:OnDelete:CASCADE"`
}

type LessonFile struct {
	UUID       uuid.UUID `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	Title      string    `json:"title"`
	FileURL    string    `json:"file_url"`
	LessonUUID uuid.UUID `json:"lesson_uuid" gorm:"type:uuid;index"`
}

func NewLesson(title, description string, difficultyID, order int, courseUUID uuid.UUID) *Lesson {
	return &Lesson{
		UUID:         uuid.New(),
		Title:        title,
		Description:  description,
		DifficultyID: difficultyID,
		Order:        order,
		CourseUUID:   courseUUID,
		Exercises:    []Exercise{},
	}
}
