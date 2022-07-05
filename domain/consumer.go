package domain

type Consumer interface {
	Subscribe(f func(message *Message) error)
	Close()
}
