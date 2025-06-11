package inbound

type MessageProcessor func(messageBody []byte, timeoutSeconds int) error

type Consumer interface {
	ProcessMessage(queue string, timeoutSeconds int, maxAttempts int, messageBody []byte) error
}
