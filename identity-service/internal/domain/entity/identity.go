package entity

import (
	"time"

	"github.com/google/uuid"
)

type Identity struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	UserUUID         uuid.UUID `json:"user_uuid" gorm:"unique"`
	Provider         string    `json:"provider"`
	Role             string    `json:"role"`
	Email            string    `json:"email"`
	PasswordHash     string    `json:"-"`
	IsConfirmEmail   bool      `json:"is_confirm_email"`
	VerificationCode string    `json:"verification_code"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func NewIdentity(email, passwordHash string) (*Identity, error) {
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}

	return &Identity{
		UserUUID:       uuid.New(),
		Provider:       "local",
		Email:          email,
		PasswordHash:   passwordHash,
		IsConfirmEmail: false,
		Role:           "user",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (i *Identity) UpdateEmail(email string) error {
	if err := ValidateEmail(email); err != nil {
		return err
	}

	i.Email = email
	i.UpdatedAt = time.Now()
	return nil
}

func (i *Identity) UpdatePassword(passwordHash string) {
	i.PasswordHash = passwordHash
	i.UpdatedAt = time.Now()
}

func (i *Identity) AddVerificationCode(code string) {
	i.VerificationCode = code
}

func (i *Identity) RemoveVerificationCode() {
	i.VerificationCode = ""
}

func (i *Identity) ConfirmEmail() {
	i.IsConfirmEmail = true
}

func (i *Identity) IsConfirmedEmail() bool {
	if i.IsConfirmEmail {
		return true
	}
	return false
}
