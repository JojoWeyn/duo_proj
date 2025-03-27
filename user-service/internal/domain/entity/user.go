package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID         uuid.UUID     `json:"uuid" gorm:"unique;type:uuid;primaryKey"`
	Login        string        `json:"login"`
	Name         string        `json:"name"`
	SecondName   string        `json:"second_name"`
	LastName     string        `json:"last_name"`
	RankID       int           `json:"rank_id"`
	Rank         Rank          `json:"rank" gorm:"foreignKey:RankID;constraint:OnDelete:SET NULL"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Avatar       string        `json:"avatar"`
	Achievements []Achievement `json:"achievements" gorm:"many2many:user_achievements;"`
}

type Leaderboard struct {
	UserUUID    uuid.UUID `json:"user_uuid"`
	Login       string    `json:"login"`
	Name        string    `json:"name"`
	SecondName  string    `json:"second_name"`
	LastName    string    `json:"last_name"`
	TotalPoints int       `json:"total_points"`
	Avatar      string    `json:"avatar"`
	Rank        int       `json:"rank"`
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
func (u *User) Validate() error {
	if u.Login == "" {
		return errors.New("login is empty")
	}

	if len(u.Login) >= 20 {
		return errors.New("login is too long")
	}

	if u.Name == "" {
		return errors.New("name is empty")
	}

	if len(u.Name) >= 30 {
		return errors.New("name is too long")
	}

	if u.SecondName == "" {
		return errors.New("second name is empty")
	}

	if len(u.SecondName) >= 30 {
		return errors.New("second name is too long")
	}

	if len(u.LastName) >= 30 {
		return errors.New("last name is too long")
	}

	return nil
}
