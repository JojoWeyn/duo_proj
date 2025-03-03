package service

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"time"
)

type VerificationService struct {
	// Add necessary fields, e.g., SMTP server details
}

func NewVerificationService() *VerificationService {
	return &VerificationService{}
}

func (vs *VerificationService) GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (vs *VerificationService) SendVerificationCode(email, code string) error {
	auth := smtp.PlainAuth("", "your-email@example.com", "your-password", "smtp.example.com")
	to := []string{email}
	msg := []byte("Subject: Verification Code\n\nYour verification code is: " + code)
	return smtp.SendMail("smtp.example.com:587", auth, "your-email@example.com", to, msg)
}
