package common

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"github.com/hashicorp/go-hclog"
	ampq "github.com/rabbitmq/amqp091-go"
)

var l = hclog.Default().Named("MessagingClient")

type IMessagingClient interface {
	Subscribe(SubscribeOpt)
	Publish(exchangeName, routingKey string, body any) error
}

type MessagingClient struct {
	conn *ampq.Connection
}

func failOnError(err error, msg string) {
	if err != nil {
		l.Error(msg, "error", err)
	}
}

func NewMessagingClient() *MessagingClient {
	RABBIT_URL := os.Getenv("RABBIT_URL")
	if RABBIT_URL == "" {
		panic("RabbitMQ URL is not set")
	}

	c, err := ampq.Dial(RABBIT_URL)
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
	}
	return &MessagingClient{c}
}

type SubscribeOpt struct {
	ExchangeName string
	RoutingKeys  []string
	QueueName    string
	HandlerFunc  func(ampq.Delivery)
}

func (client *MessagingClient) Publish(exchangeName, routingKey string, body any) error {
	json, err := json.Marshal(body)
	if err != nil {
		return err
	}

	ch, err := client.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = ch.PublishWithContext(
		ctx,
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		ampq.Publishing{
			ContentType: "application/json",
			Body:        json,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}
	return nil
}

func (m *MessagingClient) Subscribe(opt SubscribeOpt) {
	ch, err := m.conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
	}

	err = ch.ExchangeDeclare(
		opt.ExchangeName, // name
		"topic",          // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,
	)
	if err != nil {
		failOnError(err, "Failed to declare an exchange")
	}

	q, err := ch.QueueDeclare(
		opt.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		failOnError(err, "Failed to declare a queue")
	}

	for _, routingKey := range opt.RoutingKeys {
		err := ch.QueueBind(
			q.Name,           // queue name
			routingKey,       // routing key
			opt.ExchangeName, // exchange
			false,            // no-wait
			nil,
		)
		failOnError(err, fmt.Sprintf("Failed to bind a queue: %s, routing: %s", q.Name, routingKey))

	}
	if err != nil {
		failOnError(err, "Failed to declare a queue")
	}

	msg, err := ch.Consume(
		q.Name, // name,
		"",     // consumer
		false,  // auto-ack
		true,   // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		failOnError(err, "Failed to register a consumer")
	}

	go subscribeToMessages(msg, opt.HandlerFunc)
	l.Info(fmt.Sprintf("Subscribed to %s", opt.ExchangeName))

}

func subscribeToMessages(msg <-chan ampq.Delivery, handler func(ampq.Delivery)) {
	for d := range msg {
		go handler(d)
	}
}
