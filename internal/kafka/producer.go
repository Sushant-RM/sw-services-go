package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer      *kafkago.Writer
	createTopic string
	updateTopic string
}

func NewProducer() *Producer {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}

	createTopic := os.Getenv("KAFKA_TOPIC_CREATE_CONNECTION")
	if createTopic == "" {
		createTopic = "save-sw-connection"
	}

	updateTopic := os.Getenv("KAFKA_TOPIC_UPDATE_CONNECTION")
	if updateTopic == "" {
		updateTopic = "update-sw-connection"
	}

	writer := &kafkago.Writer{
		Addr:     kafkago.TCP(broker),
		Balancer: &kafkago.LeastBytes{},
	}

	return &Producer{
		writer:      writer,
		createTopic: createTopic,
		updateTopic: updateTopic,
	}
}

func (p *Producer) PublishEvent(topic, key string, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	err = p.writer.WriteMessages(
		ctx,
		kafkago.Message{
			Topic: topic,
			Key:   []byte(key),
			Value: payload,
		},
	)

	if err != nil {
		log.Printf("Kafka publish failed to topic %s: %v", topic, err)
		return err
	}

	log.Printf("Kafka event successfully published to topic %s | Key=%s", topic, key)
	return nil
}

func (p *Producer) PublishConnectionCreated(event interface{}, key string) error {
	return p.PublishEvent(p.createTopic, key, event)
}

func (p *Producer) PublishConnectionUpdated(event interface{}, key string) error {
	return p.PublishEvent(p.updateTopic, key, event)
}
