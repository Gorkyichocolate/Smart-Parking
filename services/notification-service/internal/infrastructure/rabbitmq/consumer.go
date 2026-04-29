package rabbitmq

import (
	"encoding/json"
	"log"

	"smart-parking/pkg/rabbitmq/consumer"
)

type PaymentCreated struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

func StartPaymentConsumer(cons *consumer.Consumer) {
	handler := func(body []byte) {
		var event PaymentCreated
		json.Unmarshal(body, &event)

		log.Printf("Send notification for payment: %+v\n", event)
	}

	cons.Consume(
		"notification.queue",
		"events",
		"payment.created",
		handler,
	)
}
