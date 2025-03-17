package entity

import (
	"github.com/google/uuid"
	"time"
)

type UserExerciseProgress struct {
	UUID         uuid.UUID `json:"uuid" gorm:"type:uuid;primaryKey;"`
	UserUUID     uuid.UUID `json:"user_uuid" gorm:"type:uuid;index"`
	ExerciseUUID uuid.UUID `json:"exercise_uuid" gorm:"type:uuid;index"`
	Completed    bool      `json:"completed"`
	TotalPoints  int       `json:"total_points"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func NewUserExerciseProgress(userUUID, exerciseUUID uuid.UUID, points int, completed bool) UserExerciseProgress {
	return UserExerciseProgress{
		UUID:         uuid.New(),
		UserUUID:     userUUID,
		ExerciseUUID: exerciseUUID,
		Completed:    completed,
		TotalPoints:  points,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
