package entity

import "github.com/google/uuid"

type Lesson struct {
	UUID         string     `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	Title        string     `json:"title"`
	Description  string     `json:"description" gorm:"type:text"`
	DifficultyID int        `json:"difficulty_id" gorm:"index"`
	Difficulty   Difficulty `json:"difficulty" gorm:"foreignKey:DifficultyID;references:ID"`
	Order        int        `json:"order"`
	CourseUUID   uuid.UUID  `json:"course_uuid" gorm:"type:uuid;index"`
	Course       Course     `json:"course" gorm:"foreignKey:CourseUUID;references:UUID"`
	Exercises    []Exercise `json:"exercises" gorm:"foreignKey:LessonUUID;constraint:OnDelete:CASCADE"`
}
