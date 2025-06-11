package rabbitmq

import (
	"context"
	"fmt"
	"frog-go/internal/core/ports/outbound/messagebus"
	"frog-go/internal/utils/logger"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	log     *logger.Logger
	conn    *amqp.Connection
	channel *amqp.Channel
	amqpURI string
}

type RabbitMessage struct {
	delivery amqp.Delivery
}

func (m *RabbitMessage) Body() []byte {
	return m.delivery.Body
}

func (m *RabbitMessage) Ack() error {
	return m.delivery.Ack(false)
}

type rabbitConsumer struct {
	ch         *amqp.Channel
	deliveries <-chan amqp.Delivery
	msgChan    chan messagebus.Message
	log        *logger.Logger
}

func (c *rabbitConsumer) Messages() <-chan messagebus.Message {
	return c.msgChan
}

func (c *rabbitConsumer) Close() error {
	if err := c.ch.Close(); err != nil {
		c.log.Error("Failed to close consumer channel: %v", err)
		return err
	}
	return nil
}

func NewRabbitMQ(user, password, host, port string) (messagebus.MessageBus, error) {
	log := logger.NewLogger("RabbitMQ")

	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port)

	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	log.Start("Host: %s:%s | User: %s", host, port, user)

	return &RabbitMQ{
		log:     log,
		conn:    conn,
		channel: ch,
		amqpURI: amqpURI,
	}, nil
}

func (r *RabbitMQ) reconnect() error {
	r.log.Warn("Attempting to reconnect to RabbitMQ...")

	conn, err := amqp.Dial(r.amqpURI)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	if r.conn != nil {
		r.conn.Close()
	}
	if r.channel != nil {
		r.channel.Close()
	}

	r.conn = conn
	r.channel = ch
	r.log.Info("RabbitMQ successfully reconnected.")
	return nil
}

func (r *RabbitMQ) ensureQueueExists(queueName string) error {
	if queueName == "" {
		return fmt.Errorf("queue name cannot be empty")
	}

	_, err := r.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		r.log.Error("Failed to declare queue '%s': %v", queueName, err)
	}
	return err
}

func (r *RabbitMQ) SendMessage(queueName string, body []byte) error {
	if r.channel.IsClosed() {
		if err := r.reconnect(); err != nil {
			return err
		}
	}

	if err := r.ensureQueueExists(queueName); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.channel.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		r.log.Error("Failed to send message to queue '%s': %v\nPayload: %s", queueName, err, string(body))
		return err
	}

	r.log.Info("Message sent to queue '%s': %s", queueName, string(body))
	return nil
}

func (r *RabbitMQ) Consume(queueName string) (messagebus.Consumer, error) {
	ch, err := r.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	if err := r.ensureQueueExists(queueName); err != nil {
		return nil, err
	}

	if err := ch.Qos(3, 0, false); err != nil {
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	deliveries, err := ch.Consume(
		queueName,
		"",
		false, // autoAck
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.log.Error("Failed to consume from queue '%s': %v", queueName, err)
		return nil, err
	}

	msgChan := make(chan messagebus.Message)

	consumer := &rabbitConsumer{
		ch:         ch,
		deliveries: deliveries,
		msgChan:    msgChan,
		log:        r.log,
	}

	go func() {
		for d := range deliveries {
			msgChan <- &RabbitMessage{delivery: d}
		}
		close(msgChan)
	}()

	return consumer, nil
}

func (r *RabbitMQ) DeleteQueue(queueName string) error {
	if r.channel.IsClosed() {
		if err := r.reconnect(); err != nil {
			return err
		}
	}
	_, err := r.channel.QueueDelete(queueName, false, false, false)
	return err
}

func (r *RabbitMQ) Close() {
	if err := r.channel.Close(); err != nil {
		r.log.Error("Failed to close channel: %v", err)
	}
	if err := r.conn.Close(); err != nil {
		r.log.Error("Failed to close connection: %v", err)
	}
}
