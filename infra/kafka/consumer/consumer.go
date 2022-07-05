package consumer

import (
	"br.com.thyagoluciano.poc/domain"
	"br.com.thyagoluciano.poc/infra/kafka/types"
	"errors"
	"fmt"
)

// ErrUndefinedProvider : undefined provider error.
var ErrUndefinedProvider = errors.New("undefinedProvider")

type Subscriber struct {
	provider types.Provider
	consumer domain.Consumer
}

func NewConsumer(kafkaURLs []string, topic string, groupID string, provider types.Provider) (domain.Consumer, error) {
	switch provider {
	case types.Sarama:
		saramaConsumer, err := NewSaramaConsumer(kafkaURLs, topic, groupID)
		if err != nil {
			return nil, fmt.Errorf("error: %w, topic: %s, groupID: %s", err, topic, groupID)
		}

		return &Subscriber{
			provider: provider,
			consumer: saramaConsumer,
		}, nil
	case types.Segmentio:
		segmentioConsumer, err := NewSegmentioConsumer(kafkaURLs, topic, groupID)
		if err != nil {
			return nil, fmt.Errorf("error: %w, topic: %s, groupID: %s", err, topic, groupID)
		}

		return &Subscriber{
			provider: provider,
			consumer: segmentioConsumer,
		}, nil
	case types.Confluent:
		confluentConsumer, err := NewConfluentConsumer(kafkaURLs, topic, groupID)
		if err != nil {
			return nil, fmt.Errorf("error: %w, topic: %s, groupID: %s", err, topic, groupID)
		}

		return &Subscriber{
			provider: provider,
			consumer: confluentConsumer,
		}, nil
	}

	return nil, fmt.Errorf("error: %w, message: %s, topic: %s, groupID: %s",
		ErrUndefinedProvider,
		"undefined provider or not informed",
		topic,
		groupID,
	)

}

// Subscribe reads and return the message from the Topic.
func (s Subscriber) Subscribe(f func(message *domain.Message) error) {
	s.consumer.Subscribe(f)
}

// Close finish consumer and stop read messages.
func (s Subscriber) Close() {
	s.consumer.Close()
}
