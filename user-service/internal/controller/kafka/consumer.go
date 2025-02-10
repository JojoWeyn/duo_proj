package kafka

import (
	"context"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type SaramaConsumer struct {
	brokers []string
	topic   string
}

func NewSaramaConsumer(brokers []string, topic string) *SaramaConsumer {
	return &SaramaConsumer{
		brokers: brokers,
		topic:   topic,
	}
}

func (c *SaramaConsumer) Start(ctx context.Context) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	client, err := sarama.NewClient(c.brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka client: %v", err)
	}
	defer client.Close()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(c.topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case message := <-partitionConsumer.Messages():
				log.Printf("Received message: %s", string(message.Value))
			case <-ctx.Done():
				log.Println("Shutting down consumer...")
				return
			}
		}
	}()

	wg.Wait()
}
