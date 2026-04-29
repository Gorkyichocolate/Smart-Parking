package main

import (
	"time"

	"payment-service/internal/infrastructer/rabbitmq"
	"smart-parking/pkg/rabbitmq/connection"
	"smart-parking/pkg/rabbitmq/publisher"
)

func main() {
	conn := connection.New("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()

	pub := publisher.New(ch, "events")
	paymentPub := rabbitmq.NewPaymentPublisher(pub)

	for {
		paymentPub.PublishPaymentCreated(rabbitmq.PaymentCreated{
			ID:     "123",
			Amount: 50,
		})

		time.Sleep(5 * time.Second)
	}
}
