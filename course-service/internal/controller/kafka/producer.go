package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"time"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

type UserProgressEvent struct {
	UserUUID     uuid.UUID `json:"user_uuid"`
	ExerciseUUID uuid.UUID `json:"exercise_uuid"`
	Points       int       `json:"points"`
	IsCorrect    bool      `json:"is_correct"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserAttemptEvent struct {
	UserUUID     uuid.UUID `json:"user_uuid"`
	QuestionUUID uuid.UUID `json:"question_uuid"`
	IsCorrect    bool      `json:"is_correct"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewProducer(brokers, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{brokers}, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (p *Producer) SendUserAttemptEvent(userUUID, questionUUID uuid.UUID, isCorrect bool) error {
	event := UserAttemptEvent{
		UserUUID:     userUUID,
		QuestionUUID: questionUUID,
		IsCorrect:    isCorrect,
		CreatedAt:    time.Now(),
	}

	msgBytes, _ := json.Marshal(event)
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msgBytes),
	}

	_, _, err := p.producer.SendMessage(msg)
	return err
}

func (p *Producer) SendUserProgressEvent(event UserProgressEvent) error {
	msgBytes, _ := json.Marshal(event)
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msgBytes),
	}

	_, _, err := p.producer.SendMessage(msg)
	return err
}

func (p *Producer) Close() {
	p.producer.Close()
}
