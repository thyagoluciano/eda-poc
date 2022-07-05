package domain

type Producer interface {
	Publish(key string, event Event) error
	Close()
}
