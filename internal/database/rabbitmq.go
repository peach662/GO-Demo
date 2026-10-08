package database

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func OpenRabbitMQ(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	return conn, nil
}

func CloseRabbitMQ(conn *amqp.Connection) error {
	return conn.Close()
}

func DeclareQueue(conn *amqp.Connection, queue string) error {

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	return nil
}

func PublishMessage(conn *amqp.Connection, queue string, message []byte) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()
	err = ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

func ConsumeMessage(conn *amqp.Connection, queue string) ([]byte, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)

	}
	defer ch.Close()
	msgs, err := ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume message: %w", err)
	}
	msg := <-msgs

	err = msg.Ack(false)
	if err != nil {
		return nil, fmt.Errorf("failed to ack message: %w", err)
	}
	return msg.Body, err

}

func StartConsumer(conn *amqp.Connection, queue string) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	msgs, err := ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		return fmt.Errorf("failed to consume message: %w", err)
	}

	go func() {
		defer ch.Close()
		for msg := range msgs {
			fmt.Printf("Received message: %s\n", msg.Body)
			if err := msg.Ack(false); err != nil {
				fmt.Printf("failed to ack message: %v\n", err)
			}
		}
	}()

	return nil
}