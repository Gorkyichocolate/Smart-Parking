package main

import (
	"notification-service/internal/infrastructure/rabbitmq"
	"smart-parking/pkg/rabbitmq/connection"
	"smart-parking/pkg/rabbitmq/consumer"
)

func main() {
	conn := connection.New("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()

	cons := consumer.New(ch)

	rabbitmq.StartPaymentConsumer(cons)
}
