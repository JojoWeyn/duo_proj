package dto

import (
	"time"
)

// UserEventDTO представляет событие пользователя для Kafka
type UserEventDTO struct {
	UUID      string    `json:"uuid"`
	Login     string    `json:"login"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

// ProgressEventDTO представляет событие прогресса для Kafka
type ProgressEventDTO struct {
	UserUUID   string    `json:"user_uuid"`
	EntityType string    `json:"entity_type"`
	EntityUUID string    `json:"entity_uuid"`
	Points     int       `json:"points"`
	Timestamp  time.Time `json:"timestamp"`
}

// QuestionAttemptEventDTO представляет событие попытки ответа на вопрос для Kafka
type QuestionAttemptEventDTO struct {
	UserUUID     string    `json:"user_uuid"`
	QuestionUUID string    `json:"question_uuid"`
	IsCorrect    bool      `json:"is_correct"`
	Timestamp    time.Time `json:"timestamp"`
}
