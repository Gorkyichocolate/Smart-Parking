package consumer

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	ch *amqp.Channel
}

func New(ch *amqp.Channel) *Consumer {
	return &Consumer{ch: ch}
}

func (c *Consumer) Consume(queueName, exchange, bindingKey string, handler func([]byte)) {
	_, err := c.ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("queue declare error: %v", err)
	}

	err = c.ch.QueueBind(
		queueName,
		bindingKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("queue bind error: %v", err)
	}

	msgs, err := c.ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("consume error: %v", err)
	}

	log.Println("waiting for messages...")

	for msg := range msgs {
		handler(msg.Body)
	}
}
