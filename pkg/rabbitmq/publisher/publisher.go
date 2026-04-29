package publisher

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch       *amqp.Channel
	exchange string
}

func New(ch *amqp.Channel, exchange string) *Publisher {
	err := ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("exchange declare error: %v", err)
	}

	return &Publisher{
		ch:       ch,
		exchange: exchange,
	}
}

func (p *Publisher) Publish(routingKey string, body []byte) {
	err := p.ch.PublishWithContext(
		context.Background(),
		p.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("publish error: %v", err)
	}
}
