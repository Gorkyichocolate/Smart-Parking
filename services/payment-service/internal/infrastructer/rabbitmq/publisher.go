package rabbitmq

import (
	"encoding/json"
	"payment-service/internal/domain"

	"smart-parking/pkg/rabbitmq/publisher"
)

type PaymentPublisher struct {
	pub *publisher.Publisher
}

func NewPaymentPublisher(pub *publisher.Publisher) *PaymentPublisher {
	return &PaymentPublisher{pub: pub}
}

func (p *PaymentPublisher) PublishPaymentCreated(payment domain.Payment) {

	body, _ := json.Marshal(payment)

	p.pub.Publish("payment.created", body)
}
