package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

type UserCreatedEvent struct {
	UUID  string `json:"uuid"`
	Login string `json:"login"`
}

type UserLoginEvent struct {
	UUID   string `json:"uuid"`
	Login  string `json:"login"`
	Action string `json:"action"`
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

func (p *Producer) SendUserCreated(uuid, login string) error {
	event := UserCreatedEvent{
		UUID:  uuid,
		Login: login,
	}

	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(value),
	}

	_, _, err = p.producer.SendMessage(msg)
	return err
}

func (p *Producer) SendUserLogin(uuid, login string) error {
	event := UserLoginEvent{
		UUID:   uuid,
		Login:  login,
		Action: "login",
	}

	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(value),
	}

	_, _, err = p.producer.SendMessage(msg)
	return err
}

func (p *Producer) Close() {
	p.producer.Close()
}
