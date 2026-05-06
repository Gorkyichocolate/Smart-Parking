package usecase

import (
	"payment-service/internal/domain"
	"payment-service/internal/repository"
)

type PaymentPublisher interface {
	PublishPaymentCreated(payment domain.Payment)
}

type PaymentUsecase struct {
	paymentRepo repository.PaymentRepository
	invoiceRepo repository.InvoiceRepository
	publisher   PaymentPublisher
}

func NewPaymentUsecase(
	p repository.PaymentRepository,
	i repository.InvoiceRepository,
	pub PaymentPublisher,
) *PaymentUsecase {
	return &PaymentUsecase{
		paymentRepo: p,
		invoiceRepo: i,
		publisher:   pub,
	}
}

func (u *PaymentUsecase) CreatePayment(p domain.Payment) error {
	// 1. payment
	if err := u.paymentRepo.CreatePayment(p); err != nil {
		return err
	}

	// 2. invoice
	invoice := domain.Invoice{
		ID:        "inv_" + p.ID,
		PaymentID: p.ID,
		UserID:    p.UserID,
		Amount:    p.Amount,
		PDFUrl:    "",
	}

	if err := u.invoiceRepo.CreateInvoice(invoice); err != nil {
		return err
	}

	// 3. event
	u.publisher.PublishPaymentCreated(p)

	return nil
}
