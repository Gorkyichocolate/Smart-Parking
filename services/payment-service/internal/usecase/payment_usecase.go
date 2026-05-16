// internal/usecase/payment_usecase.go
package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/Gorkyichocolate/smart-parking/services/payment-service/internal/domain"
	"github.com/Gorkyichocolate/smart-parking/services/payment-service/internal/repository"
)

type PaymentUseCase struct {
	paymentRepo     repository.PaymentRepository
	invoiceRepo     repository.InvoiceRepository
	rabbitPublisher interface {
		PublishNotification(userEmail, bookingID, amount string) error
	}
}

func NewPaymentUseCase(
	paymentRepo repository.PaymentRepository,
	invoiceRepo repository.InvoiceRepository,
	rabbitPublisher interface {
		PublishNotification(userEmail, bookingID, amount string) error
	},
) *PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo:     paymentRepo,
		invoiceRepo:     invoiceRepo,
		rabbitPublisher: rabbitPublisher,
	}
}

type CreatePaymentInput struct {
	BookingID     string
	UserID        string
	Amount        float64
	PaymentMethod domain.PaymentMethod
	UserEmail     string
}

func (uc *PaymentUseCase) CreatePayment(ctx context.Context, input CreatePaymentInput) (*domain.Payment, error) {
	// Валидация
	if input.BookingID == "" {
		return nil, domain.ErrInvalidPaymentData
	}
	if input.UserID == "" {
		return nil, domain.ErrInvalidPaymentData
	}
	if input.Amount <= 0 {
		return nil, domain.ErrInvalidPaymentData
	}
	if input.PaymentMethod == "" {
		return nil, domain.ErrInvalidPaymentData
	}

	// Создаем платеж
	payment := &domain.Payment{
		ID:            uuid.New().String(),
		BookingID:     input.BookingID,
		UserID:        input.UserID,
		Amount:        input.Amount,
		Status:        domain.PaymentStatusPending,
		PaymentMethod: input.PaymentMethod,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Сохраняем в БД
	if err := uc.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Асинхронная обработка платежа
	go uc.processPayment(payment, input.UserEmail)

	return payment, nil
}

func (uc *PaymentUseCase) processPayment(payment *domain.Payment, userEmail string) {
	// Симуляция обработки платежа
	time.Sleep(2 * time.Second)

	ctx := context.Background()

	// Обновляем статус платежа
	payment.Status = domain.PaymentStatusCompleted
	if err := uc.paymentRepo.UpdateStatus(ctx, payment.ID, payment.Status); err != nil {
		log.Printf("Failed to update payment status: %v", err)
		return
	}

	// Создаем инвойс
	invoice := &domain.Invoice{
		ID:        uuid.New().String(),
		PaymentID: payment.ID,
		UserID:    payment.UserID,
		Amount:    payment.Amount,
		PDFURL:    fmt.Sprintf("/invoices/%s.pdf", payment.ID),
		IssuedAt:  time.Now(),
	}

	if err := uc.invoiceRepo.Create(ctx, invoice); err != nil {
		log.Printf("Failed to create invoice: %v", err)
	} else {
		log.Printf("Invoice created: %s", invoice.ID)
	}

	// Отправляем уведомление
	if err := uc.rabbitPublisher.PublishNotification(
		userEmail,
		payment.BookingID,
		fmt.Sprintf("%.2f", payment.Amount),
	); err != nil {
		log.Printf("Failed to publish notification: %v", err)
	}

	log.Printf("Payment processed: ID=%s, BookingID=%s, Amount=%.2f",
		payment.ID, payment.BookingID, payment.Amount)
}

func (uc *PaymentUseCase) GetPayment(ctx context.Context, id string) (*domain.Payment, error) {
	return uc.paymentRepo.GetByID(ctx, id)
}
