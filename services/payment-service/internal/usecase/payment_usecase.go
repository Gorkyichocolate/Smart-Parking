package usecase

import (
	"payment-service/internal/domain"
	"payment-service/internal/repository"
)

type PaymentPublisher interface {
	PublishPaymentCreated(payment domain.Payment)
}

type PaymentUseCase struct {
	repo      repository.PaymentRepository
	publisher PaymentPublisher
}

func NewPaymentUseCase(r repository.PaymentRepository, p PaymentPublisher) *PaymentUseCase {
	return &PaymentUseCase{
		repo:      r,
		publisher: p,
	}
}

func (u *PaymentUseCase) CreatePayment(p domain.Payment) error {
	err := u.repo.CreatePayment(p)
	if err != nil {
		return err
	}

	u.publisher.PublishPaymentCreated(p)

	return nil
}
