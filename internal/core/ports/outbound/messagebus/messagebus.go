package messagebus

type Message interface {
	Body() []byte
	Ack() error
}

type Consumer interface {
	Messages() <-chan Message
	Close() error
}

type MessageBus interface {
	SendMessage(queueName string, body []byte) error
	Consume(queueName string) (Consumer, error)
	DeleteQueue(queueName string) error
	Close()
}
