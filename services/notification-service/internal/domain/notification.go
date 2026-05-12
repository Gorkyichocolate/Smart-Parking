package domain

import (
	"time"
)

type EmailNotification struct {
	ID        string
	To        string
	Subject   string
	Body      string
	IsHTML    bool
	Status    EmailStatus
	Attempts  int
	CreatedAt time.Time
	SentAt    *time.Time
	Error     string
}

type EmailStatus string

const (
	EmailStatusPending   EmailStatus = "pending"
	EmailStatusSent      EmailStatus = "sent"
	EmailStatusFailed    EmailStatus = "failed"
	EmailStatusDelivered EmailStatus = "delivered"
)

type EmailNotifier interface {
	SendEmail(notification EmailNotification) error
	ValidateEmail(email string) bool
	GetRetryDelay(attempt int) time.Duration
}
