package repository

import "payment-service/internal/domain"

type PaymentRepository interface {
	CreatePayment(payment domain.Payment) error
	GetPaymentByID(id string) (domain.Payment, error)
	DeletePayment(id string) error
}
