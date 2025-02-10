package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID       uuid.UUID `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	Name       string    `json:"name"`
	SecondName string    `json:"second_name"`
	LastName   string    `json:"last_name"`
	RankID     int       `json:"rank_id"`
	Rank       Rank      `json:"rank" gorm:"foreignKey:RankID;constraint:OnDelete:SET NULL"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Avatar     string    `json:"avatar"`
}

func NewUser(userUUID uuid.UUID) *User {
	return &User{
		UUID:      userUUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		RankID:    1,
	}
}
