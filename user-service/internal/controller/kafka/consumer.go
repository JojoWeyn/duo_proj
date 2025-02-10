package kafka

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type SaramaConsumerGroup struct {
	brokers []string
	topic   string
	groupID string
}

func NewSaramaConsumerGroup(brokers []string, topic, groupID string) *SaramaConsumerGroup {
	return &SaramaConsumerGroup{
		brokers: brokers,
		topic:   topic,
		groupID: groupID,
	}
}

func (c *SaramaConsumerGroup) Start(ctx context.Context) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(c.brokers, c.groupID, config)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err)
	}
	defer consumerGroup.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	consumer := ConsumerGroupHandler{}

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

type ConsumerGroupHandler struct{}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}
	return nil
}
