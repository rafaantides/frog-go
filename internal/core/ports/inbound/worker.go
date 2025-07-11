package inbound

type Consumer interface {
	ProcessMessage(timeoutSeconds int, messageBody []byte) error
}
