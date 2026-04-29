package connection

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func New(url string) *amqp.Connection {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}
	return conn
}