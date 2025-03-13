package entity

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Achievement struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Condition   string    `json:"condition" gorm:"type:jsonb"`
	Secret      bool      `json:"secret"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserAchievementProgress struct {
	ID            int        `json:"id" gorm:"primaryKey"`
	UserUUID      uuid.UUID  `json:"user_uuid" gorm:"type:uuid;index"`
	AchievementID int        `json:"achievement_id" gorm:"index"`
	CurrentCount  int        `json:"current_count"`
	Achieved      bool       `json:"achieved"`
	AchievedAt    *time.Time `json:"achieved_at,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Condition struct {
	Count      int      `json:"count,omitempty"`
	Action     string   `json:"action,omitempty"`
	Timeframe  string   `json:"timeframe,omitempty"`
	Stat       string   `json:"stat,omitempty"`
	TopPercent int      `json:"top_percent,omitempty"`
	ActionSeq  []string `json:"action_sequence,omitempty"`
	Secret     bool     `json:"secret,omitempty"`
}

func (a *Achievement) Validate() error {
	if a.Title == "" {
		return errors.New("title is required")
	}
	if a.Description == "" {
		return errors.New("description is required")
	}
	if a.Condition == "" {
		return errors.New("condition is required")
	}
	return nil
}
