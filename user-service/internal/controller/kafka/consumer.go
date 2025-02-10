package kafka

import (
	"log"

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

func (c *SaramaConsumer) Start() {
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

	for message := range partitionConsumer.Messages() {
		log.Printf("Received message: %s", string(message.Value))
	}
}
