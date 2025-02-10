package entity

import (
	"time"
)

type BlacklistedToken struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Token     string    `json:"token" gorm:"unique"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func NewBlacklistedToken(token string, expiresAt time.Time) *BlacklistedToken {
	return &BlacklistedToken{
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}
