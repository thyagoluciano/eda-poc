package domain

import (
	"encoding/json"
	"time"
)

type Event struct {
	EventId   string      `json:"eventId"`
	EventType string      `json:"eventType"`
	Metadata  Metadata    `json:"metadata"`
	Payload   interface{} `json:"payload,omitempty"`
}

type Context struct {
	SpanId       string      `json:"spanId"`
	TraceId      string      `json:"traceId"`
	Organization string      `json:"organization"`
	Application  string      `json:"application"`
	Channel      string      `json:"channel"`
	Custom       interface{} `json:"custom"`
}

type Metadata struct {
	Id         string    `json:"id"`
	Domain     string    `json:"domain"`
	ExternalId string    `json:"externalId"`
	CustomerId string    `json:"customerId"`
	Context    Context   `json:"context"`
	Timestamp  time.Time `json:"timestamp"`
}

func (e Event) ToEvent(message *Message) (Event, error) {
	event := Event{}
	err := json.Unmarshal(message.Value, &event)
	if err != nil {
		return Event{}, err
	}

	return event, nil
}
