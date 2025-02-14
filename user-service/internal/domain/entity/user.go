package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID       uuid.UUID `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	Login      string    `json:"login"`
	Name       string    `json:"name"`
	SecondName string    `json:"second_name"`
	LastName   string    `json:"last_name"`
	RankID     int       `json:"rank_id"`
	Rank       Rank      `json:"rank" gorm:"foreignKey:RankID;constraint:OnDelete:SET NULL"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Avatar     string    `json:"avatar"`
}

func NewUser(userUUID uuid.UUID, login string) *User {
	return &User{
		UUID:      userUUID,
		Login:     login,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		RankID:    1,
		Avatar:    "https://ybis.ru/wp-content/uploads/2023/09/solntse-kartinka-1.webp",
	}
}
