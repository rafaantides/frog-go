package inbound

type MessageProcessor func(timeoutSeconds int, messageBody []byte) error

type Consumer interface {
	ProcessMessage(timeoutSeconds int, messageBody []byte) error
}
