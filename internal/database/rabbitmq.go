package database

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"context"
	"time"
)

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
	defer func() { _ = ch.Close() }()
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
	defer func() { _ = ch.Close() }()
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
	defer func() { _ = ch.Close() }()
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

func StartConsumer(conn *amqp.Connection, queue string,redisClient *redis.Client) error {
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
		defer func() { _ = ch.Close() }()

		seen := make(map[string]struct{})
		for msg := range msgs {
			fmt.Printf("Received message: %s\n", msg.Body)

			var  eventTodo struct {
				EventType string `json:"event_type"`
				Todo struct {
				ID    int    `json:"id"`
				Title string `json:"title"`
				} `json:"todo"`
			
			}
			err := json.Unmarshal(msg.Body, &eventTodo)
			if err != nil {
				fmt.Printf("failed to unmarshal message: %v\n", err)
				fmt.Printf("message: %s\n", msg.Body)
				_ = msg.Nack(false, false)
				continue
			}
			key := fmt.Sprintf("todo:event:%s:%d", eventTodo.EventType, eventTodo.Todo.ID)
			if redisClient == nil {
				if _, ok := seen[key]; ok {
					fmt.Printf("Duplicate message: %s\n", msg.Body)
					_ = msg.Ack(false)
					continue
				}
				seen[key] = struct{}{}
				fmt.Printf("EventType: %s, Todo: %+v\n", eventTodo.EventType, eventTodo.Todo)
				_ = msg.Ack(false)
				continue
			}

			
			ctx := context.Background()
			n, err := redisClient.Exists(ctx, key).Result()
			if err != nil {
				fmt.Printf("failed to check redis key: %v\n", err)
				_ = msg.Nack(false, true)
				continue
			}
			if n > 0 {
				fmt.Printf("Duplicate message: %s\n", msg.Body)
				_ = msg.Ack(false)
				continue
			}

			fmt.Printf("EventType: %s, Todo: %+v\n", eventTodo.EventType, eventTodo.Todo)
			if err := redisClient.Set(ctx, key, eventTodo.EventType, 24*time.Hour).Err(); err != nil {
				fmt.Printf("failed to set redis key: %v\n", err)
				_ = msg.Nack(false, true)
				continue
			}
			_ = msg.Ack(false)
		}
	}()

	return nil
}