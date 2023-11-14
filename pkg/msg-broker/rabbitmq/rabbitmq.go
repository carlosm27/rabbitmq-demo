package rabbitmq

import (
	"context"
	"fmt"

	msgbroker "github.com/carlosm27/go-rabbitmq-demo/pkg/msg-broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqBroker struct {
	ch *amqp.Channel
}

func MustGetNewRabbitMqBroker(url string) (*RabbitMqBroker, func()) {
	// --- (1) ----
	// Try to create the connection
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(fmt.Errorf("failed to connect to RabbitMQ: %v", err))
	}
	// --- (2) ----
	// Try to create the channel
	ch, err := conn.Channel()
	if err != nil {
		panic(fmt.Errorf("failed to create RabbitMQ channel: %v", err))
	}
	return &RabbitMqBroker{
			ch: ch,
		}, func() {
			ch.Close()
			conn.Close()
		}
}

func (br *RabbitMqBroker) PublishMessage(ctx context.Context, queueSettings msgbroker.QueueSettings, msgBody []byte) error {
	// --- (1) ----
	// First we need to cleare the message queue where the message will be published
	q, err := br.ch.QueueDeclare(
		queueSettings.Name,
		queueSettings.Durable,
		queueSettings.AutoDelete,
		queueSettings.Exclusive,
		queueSettings.NoWait,
		queueSettings.Args,
	)
	if err != nil {
		return err
	}

	// --- (2) ----
	// If the queue is successfully created, we can procede to publish the message
	err = br.ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msgBody,
		})
	return err
}
