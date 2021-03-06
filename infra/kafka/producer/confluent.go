package producer

import (
	"br.com.thyagoluciano.poc/domain"
	"encoding/json"
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"strings"
	"time"
)

type ConfluentProducer struct {
	kafkaURLs []string
	topic     string
	producer  *kafka.Producer
}

// NewConfluentProducer create a new Confluent Producer.
func NewConfluentProducer(kafkaURLs []string, topic string) (*ConfluentProducer, error) {
	prc, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": strings.Join(kafkaURLs, ",")})
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	go func() {
		for e := range prc.Events() {
			switch ev := e.(type) { //nolint:gocritic
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &ConfluentProducer{
		kafkaURLs: kafkaURLs,
		topic:     topic,
		producer:  prc,
	}, nil
}

func (confluent ConfluentProducer) Publish(key string, event domain.Event) error {
	message, err := json.Marshal(event)

	if err != nil {
		//logger.Log("parse", "failed to encode json event")
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &confluent.topic,
			Partition: kafka.PartitionAny,
			Offset:    0,
			Metadata:  nil,
			Error:     nil,
		},
		Value:         message,
		Key:           []byte(key),
		Timestamp:     time.Time{},
		TimestampType: 0,
		Opaque:        nil,
		Headers:       nil,
	}
	err = confluent.producer.Produce(msg, nil)
	confluent.producer.Flush(0)

	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	return nil
}

func (confluent ConfluentProducer) Close() {
	confluent.producer.Close()
}
