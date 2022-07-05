package producer

import (
	"br.com.thyagoluciano.poc/domain"
	"br.com.thyagoluciano.poc/infra/kafka/types"
	"fmt"
)

func New(kafkaURLs []string, topic string, provider types.Provider) (domain.Producer, error) {
	switch provider {
	case types.Sarama:
		saramaProducer, err := NewSaramaProducer(kafkaURLs, topic)
		if err != nil {
			return nil, fmt.Errorf("error: %w, topic: %s", err, topic)
		}

		return saramaProducer, nil

	case types.Segmentio:
		segmentioProducer, err := NewSegmentioProducer(kafkaURLs, topic)
		if err != nil {
			return nil, fmt.Errorf("error: %w, topic: %s", err, topic)
		}

		return segmentioProducer, err

	case types.Confluent:
		confluentProducer, err := NewConfluentProducer(kafkaURLs, topic)
		if err != nil {
			return nil, fmt.Errorf("error: %w, topic: %s", err, topic)
		}

		return confluentProducer, nil
	}

	return nil, fmt.Errorf("error> %w, message: %s, topic: %s",
		types.ErrUndefinedProvider,
		"undefined provider or not informade",
		topic,
	)
}
