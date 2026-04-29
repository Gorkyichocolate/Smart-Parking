package domain

import "time"

type Invoice struct {
	ID        string
	PaymentID string
	UserID    string
	Amount    float64
	PDFUrl    string
	IssuedAt  time.Time
}
