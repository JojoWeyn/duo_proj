package dto

import (
	"encoding/json"
	"time"
)

type UserAchievementsDTO struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Secret       bool      `json:"secret"`
	CurrentCount int       `json:"current_count"`
	Achieved     bool      `json:"achieved"`
	CreatedAt    time.Time `json:"created_at"`
	AchievedAt   time.Time `json:"achieved_at"`
}

type AchievementsDTO struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Condition   json.RawMessage `json:"condition"`
	Secret      bool            `json:"secret"`
	CreatedAt   time.Time       `json:"created_at"`
}
