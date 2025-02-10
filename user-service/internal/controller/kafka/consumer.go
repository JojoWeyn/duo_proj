package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type userUseCase interface {
	CreateUser(ctx context.Context, uuid uuid.UUID) error
}

type SaramaConsumerGroup struct {
	brokers     []string
	topic       string
	groupID     string
	userUsecase userUseCase
}

func NewSaramaConsumerGroup(userUsecase userUseCase, brokers []string, topic, groupID string) *SaramaConsumerGroup {
	return &SaramaConsumerGroup{
		brokers:     brokers,
		topic:       topic,
		groupID:     groupID,
		userUsecase: userUsecase,
	}
}

func (c *SaramaConsumerGroup) Start(ctx context.Context) {
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

	consumer := ConsumerGroupHandler{useCase: c.userUsecase}

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

type ConsumerGroupHandler struct {
	useCase userUseCase
}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (c ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)

		var msg struct {
			UUID string `json:"uuid"`
		}

		if err := json.Unmarshal(message.Value, &msg); err != nil {
			log.Printf("Failed to parse message JSON: %v", err)
			continue
		}

		receivedUUID, err := uuid.Parse(msg.UUID)
		if err != nil {
			log.Printf("Invalid UUID format: %v", err)
			continue
		}
		if err := c.useCase.CreateUser(context.Background(), receivedUUID); err != nil {
			log.Printf("Error creating user: %v", err)
		}
		log.Printf("User with UUID %s successfully created", receivedUUID)
		session.MarkMessage(message, "")
	}
	return nil
}
