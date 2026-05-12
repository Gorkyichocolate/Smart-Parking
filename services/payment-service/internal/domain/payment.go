// internal/domain/payment.go
package domain

import (
	"time"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
)

type PaymentMethod string

const (
	PaymentMethodCard PaymentMethod = "card"
	PaymentMethodCash PaymentMethod = "cash"
	PaymentMethodQR   PaymentMethod = "qr"
)

type Payment struct {
	ID            string        `json:"id"`
	BookingID     string        `json:"booking_id"`
	UserID        string        `json:"user_id"`
	Amount        float64       `json:"amount"`
	Status        PaymentStatus `json:"status"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// Domain errors
var (
	ErrPaymentNotFound      = &DomainError{Code: "PAYMENT_NOT_FOUND", Message: "payment not found"}
	ErrInvalidPaymentData   = &DomainError{Code: "INVALID_PAYMENT_DATA", Message: "invalid payment data"}
	ErrPaymentAlreadyExists = &DomainError{Code: "PAYMENT_ALREADY_EXISTS", Message: "payment already exists"}
)

type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}
