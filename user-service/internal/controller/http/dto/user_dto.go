package dto

import (
	"time"

	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/google/uuid"
)

// UserDTO представляет данные пользователя для передачи между слоями
type UserDTO struct {
	UUID            uuid.UUID `json:"uuid"`
	Login           string    `json:"login"`
	Name            string    `json:"name"`
	SecondName      string    `json:"second_name"`
	LastName        string    `json:"last_name"`
	RankID          int       `json:"rank_id"`
	Rank            RankDTO   `json:"rank,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Avatar          string    `json:"avatar"`
	TotalPoints     int       `json:"total_points"`
	FinishedCourses int64     `json:"finished_courses"`
}

// RankDTO представляет данные о ранге пользователя
type RankDTO struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// UserUpdateDTO представляет данные для обновления пользователя
type UserUpdateDTO struct {
	Login      *string `json:"login,omitempty"`
	Name       *string `json:"name,omitempty"`
	SecondName *string `json:"second_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
}

// LeaderboardDTO представляет данные для таблицы лидеров
type LeaderboardDTO struct {
	UserUUID    uuid.UUID `json:"user_uuid"`
	Login       string    `json:"login"`
	Name        string    `json:"name"`
	SecondName  string    `json:"second_name"`
	LastName    string    `json:"last_name"`
	TotalPoints int       `json:"total_points"`
	Avatar      string    `json:"avatar"`
	Rank        int       `json:"rank"`
}

// UserAvatarResponseDTO представляет ответ на обновление аватара
type UserAvatarResponseDTO struct {
	AvatarURL string `json:"avatar_url"`
}

// StreakResponseDTO представляет ответ с информацией о серии дней
type StreakResponseDTO struct {
	Days int `json:"days"`
}

func ToUserDTO(u *entity.User) UserDTO {
	return UserDTO{
		UUID:            u.UUID,
		Login:           u.Login,
		Name:            u.Name,
		SecondName:      u.SecondName,
		LastName:        u.LastName,
		RankID:          u.RankID,
		Rank:            ToRankDTO(&u.Rank),
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
		Avatar:          u.Avatar,
		TotalPoints:     u.TotalPoints,
		FinishedCourses: u.FinishedCourses,
	}
}

func ToLeaderboardDTO(e entity.Leaderboard) LeaderboardDTO {
	return LeaderboardDTO{
		UserUUID:    e.UserUUID,
		Login:       e.Login,
		Name:        e.Name,
		SecondName:  e.SecondName,
		LastName:    e.LastName,
		TotalPoints: e.TotalPoints,
		Avatar:      e.Avatar,
		Rank:        e.Rank,
	}
}

func ToRankDTO(r *entity.Rank) RankDTO {
	if r == nil {
		return RankDTO{}
	}
	return RankDTO{
		ID:    r.ID,
		Title: r.Name,
	}
}
