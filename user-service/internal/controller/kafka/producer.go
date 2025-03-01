package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

type UserEvent struct {
	UUID   string `json:"uuid"`
	Login  string `json:"login"`
	Action string `json:"action"`
}

func NewProducer(brokers, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer([]string{brokers}, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (p *Producer) SendUserEvent(uuid, login, action string) error {
	event := UserEvent{
		UUID:   uuid,
		Login:  login,
		Action: action,
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
