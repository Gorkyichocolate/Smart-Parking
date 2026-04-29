package domain

import "time"

type Payment struct {
	ID            string
	BookingID     string
	UserID        string
	Amount        float64
	Status        string
	PaymentMethod string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
