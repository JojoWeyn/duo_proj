package kafka

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ProgressUseCase interface {
	CheckProgress(ctx context.Context, userID uuid.UUID, progressUUID uuid.UUID) bool
	AddProgress(ctx context.Context, userID uuid.UUID, entityType string, entityUUID uuid.UUID, points int, createdAt time.Time) error
}

type ProgressConsumer struct {
	brokers         []string
	topic           string
	groupID         string
	progressUseCase ProgressUseCase
}

func NewProgressConsumer(brokers []string, topic, groupID string, progressUseCase ProgressUseCase) *ProgressConsumer {
	return &ProgressConsumer{
		brokers:         brokers,
		topic:           topic,
		groupID:         groupID,
		progressUseCase: progressUseCase,
	}
}

func (c *ProgressConsumer) Start(ctx context.Context) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(c.brokers, c.groupID, config)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err)
	}
	defer consumerGroup.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	consumer := ProgressConsumerGroupHandler{progressUseCase: c.progressUseCase}

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{c.topic}, &consumer); err != nil {
				if ctx.Err() != nil {
					log.Println("Consumer loop exiting due to context cancellation")
					return
				}
				log.Printf("Error from consumer: %v", err)
			}
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}

}

type ProgressConsumerGroupHandler struct {
	progressUseCase ProgressUseCase
}

func (ProgressConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ProgressConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (c ProgressConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)

		var msg struct {
			UserUUID   string    `json:"user_uuid"`
			EntityType string    `json:"entity_type"`
			EntityUUID string    `json:"entity_uuid"`
			Points     int       `json:"points"`
			IsCorrect  bool      `json:"is_correct"`
			CreatedAt  time.Time `json:"created_at"`
		}

		if err := json.Unmarshal(message.Value, &msg); err != nil {
			log.Printf("Failed to parse message JSON: %v", err)
			continue
		}

		entityUUID := uuid.MustParse(msg.EntityUUID)
		userUUID := uuid.MustParse(msg.UserUUID)

		isCompleted := c.progressUseCase.CheckProgress(context.Background(), userUUID, entityUUID)

		if !isCompleted {
			if err := c.progressUseCase.AddProgress(context.Background(), userUUID, msg.EntityType, entityUUID, msg.Points, msg.CreatedAt); err != nil {
				log.Printf("Error adding progress: %v", err)
			}
		}

		session.MarkMessage(message, "")
	}
	return nil
}
