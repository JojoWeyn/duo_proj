package entity

import "github.com/google/uuid"

type Course struct {
	UUID         uuid.UUID    `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	Title        string       `json:"title"`
	Description  string       `json:"description" gorm:"type:text"`
	TypeID       int          `json:"type_id" gorm:"index"`
	CourseType   CourseType   `json:"course_type" gorm:"foreignKey:TypeID;references:ID"`
	Lessons      []Lesson     `json:"lessons" gorm:"foreignKey:CourseUUID;constraint:OnDelete:CASCADE"`
	DifficultyID int          `json:"difficulty_id" gorm:"index"`
	Difficulty   Difficulty   `json:"difficulty" gorm:"foreignKey:DifficultyID;references:ID"`
	CourseFiles  []CourseFile `json:"course_files" gorm:"foreignKey:CourseUUID;constraint:OnDelete:CASCADE"`
}

type CourseFile struct {
	UUID       uuid.UUID `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	Title      string    `json:"title"`
	FileURL    string    `json:"file_url"`
	CourseUUID uuid.UUID `json:"course_uuid" gorm:"type:uuid;index"`
}

func NewCourse(title, description string, typeID, difficultyID int) *Course {
	return &Course{
		UUID:         uuid.New(),
		Title:        title,
		Description:  description,
		TypeID:       typeID,
		DifficultyID: difficultyID,
		Lessons:      []Lesson{},
	}
}
