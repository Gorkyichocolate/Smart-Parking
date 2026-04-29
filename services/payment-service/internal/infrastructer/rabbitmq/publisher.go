package rabbitmq


import (
	"encoding/json"

	"smart-parking/pkg/rabbitmq/publisher"
)

type PaymentPublisher struct {
	pub *publisher.Publisher
}

func NewPaymentPublisher(pub *publisher.Publisher) *PaymentPublisher {
	return &PaymentPublisher{pub: pub}
}

type PaymentCreated struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

func (p *PaymentPublisher) PublishPaymentCreated(event PaymentCreated) {
	body, _ := json.Marshal(event)

	p.pub.Publish("payment.created", body)
}
