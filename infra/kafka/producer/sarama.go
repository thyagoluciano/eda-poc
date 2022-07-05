package producer

import (
	"br.com.thyagoluciano.poc/domain"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"time"
)

type SaramaProducer struct {
	kafkaURLs []string
	topic     string
	producer  sarama.SyncProducer
}

func NewSaramaProducer(kafkaURLs []string, topic string) (*SaramaProducer, error) {
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Retry.Max = 5
	saramaCfg.Producer.RequiredAcks = sarama.WaitForAll
	saramaCfg.Producer.Return.Successes = true

	prd, err := sarama.NewSyncProducer(kafkaURLs, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	return &SaramaProducer{
		kafkaURLs: kafkaURLs,
		topic:     topic,
		producer:  prd,
	}, nil
}

func (saramaProducer SaramaProducer) Publish(key string, event domain.Event) error {
	message, err := json.Marshal(event)

	if err != nil {
		//logger.Log("parse", "failed to encode json event")
	}

	msg := &sarama.ProducerMessage{
		Topic:     saramaProducer.topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.StringEncoder(message),
		Headers:   nil,
		Metadata:  nil,
		Offset:    0,
		Partition: 0,
		Timestamp: time.Time{},
	}

	_, _, err = saramaProducer.producer.SendMessage(msg)

	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	return nil
}

func (saramaProducer SaramaProducer) Close() {
	if saramaProducer.producer != nil {
		_ = saramaProducer.producer.Close()
	}
}
