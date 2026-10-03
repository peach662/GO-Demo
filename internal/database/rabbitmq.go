package database

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"fmt"
	
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch *amqp.Channel
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
