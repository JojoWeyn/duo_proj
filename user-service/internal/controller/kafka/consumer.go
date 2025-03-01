package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type userUseCase interface {
	CreateUser(ctx context.Context, uuid uuid.UUID, Login string) error
}

type achievementUseCase interface {
	CheckAchievements(ctx context.Context, userID uuid.UUID, action string) error
}

type SaramaConsumerGroup struct {
	brokers            []string
	topic              string
	groupID            string
	userUsecase        userUseCase
	achievementUsecase achievementUseCase
}

func NewSaramaConsumerGroup(userUsecase userUseCase, achievementUsecase achievementUseCase, brokers []string, topic, groupID string) *SaramaConsumerGroup {
	return &SaramaConsumerGroup{
		brokers:            brokers,
		topic:              topic,
		groupID:            groupID,
		userUsecase:        userUsecase,
		achievementUsecase: achievementUsecase,
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

	consumer := ConsumerGroupHandler{userUseCase: c.userUsecase, achievementUsecase: c.achievementUsecase}

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
	userUseCase        userUseCase
	achievementUsecase achievementUseCase
}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (c ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	if c.userUseCase == nil {
		log.Println("Error: userUseCase is nil")
		return fmt.Errorf("userUseCase is nil")
	}

	if c.achievementUsecase == nil {
		log.Println("Error: achievementUsecase is nil")
		return fmt.Errorf("achievementUsecase is nil")
	}

	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)

		var msg struct {
			UUID   string `json:"uuid"`
			Login  string `json:"login"`
			Action string `json:"action"`
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

		if msg.Action == "" {
			if err := c.userUseCase.CreateUser(context.Background(), receivedUUID, msg.Login); err != nil {
				log.Printf("Error creating user: %v", err)
			} else {
				log.Printf("User with UUID %s successfully created", receivedUUID)
			}
		} else {
			err := c.achievementUsecase.CheckAchievements(context.Background(), receivedUUID, msg.Action)
			if err != nil {
				log.Printf("Error checking achievements: %v", err)
			}
		}

		session.MarkMessage(message, "")
	}
	return nil
}
