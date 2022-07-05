package types

import "errors"

type Provider string

var ErrUndefinedProvider = errors.New("undefinedProvider")

const (
	// Sarama provider.
	Sarama Provider = "sarama"
	// Segmentio provider.
	Segmentio Provider = "segmentio"
	// Confluent provider.
	Confluent Provider = "confluent"
)

func (p Provider) String() string {
	return string(p)
}
