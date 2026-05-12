// internal/domain/invoice.go
package domain

import (
	"time"
)

type Invoice struct {
	ID        string    `json:"id"`
	PaymentID string    `json:"payment_id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	PDFURL    string    `json:"pdf_url"`
	IssuedAt  time.Time `json:"issued_at"`
}

// Invoice specific errors
var (
	ErrInvoiceNotFound = &DomainError{Code: "INVOICE_NOT_FOUND", Message: "invoice not found"}
)
