package domain

import "context"

type Service interface {
	Producer(ctx context.Context, cmd Command) error
	Consumer(message *Message) error
}
