package entity

import "github.com/google/uuid"

type Exercise struct {
	UUID        uuid.UUID  `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	Title       string     `json:"title"`
	Description string     `json:"description" gorm:"type:text"`
	Points      int        `json:"points"`
	Order       int        `json:"order"`
	LessonUUID  uuid.UUID  `json:"lesson_uuid" gorm:"type:uuid;index"`
	Lesson      Lesson     `json:"lesson" gorm:"foreignKey:LessonUUID;references:UUID"`
	Questions   []Question `json:"questions" gorm:"foreignKey:ExerciseUUID;constraint:OnDelete:CASCADE"`
}

func NewExercise(title, description string, points, order int, lessonUUID uuid.UUID) *Exercise {
	return &Exercise{
		UUID:        uuid.New(),
		Title:       title,
		Description: description,
		Points:      points,
		Order:       order,
		LessonUUID:  lessonUUID,
		Questions:   []Question{},
	}
}
