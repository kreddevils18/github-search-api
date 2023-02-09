package rabbitmq

import (
	"context"
	"fmt"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/config"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMd struct {
	channel *amqp091.Channel
}

func NewRabbitMQ(config *config.Config) (*RabbitMd, error) {
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s",
		config.RabbitMQ.User,
		config.RabbitMQ.Password,
		config.RabbitMQ.Host,
		config.RabbitMQ.Port))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMd{ch}, nil
}

func (r RabbitMd) GetChannel() *amqp091.Channel {
	return r.channel
}

func (r RabbitMd) QueueDeclare(name string) (*amqp091.Queue, error) {
	q, err := r.channel.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil)

	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (r RabbitMd) Publish(name string, data []byte) error {
	ctx := context.Background()
	err := r.channel.PublishWithContext(
		ctx,
		"",
		name,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp091.Persistent,
		})

	return err
}

func (r RabbitMd) Consume(name string) (<-chan amqp091.Delivery, error) {
	msgs, err := r.channel.Consume(
		name,
		"",
		false,
		false,
		false,
		false,
		nil)

	if err != nil {
		return nil, err
	}

	return msgs, nil
}
