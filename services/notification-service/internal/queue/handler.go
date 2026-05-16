// services/notification-service/internal/queue/handler.go
package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/Gorkyichocolate/smart-parking/services/notification-service/internal/email"

	"github.com/Gorkyichocolate/smart-parking/pkg/metrics"
)

type MessageHandler interface {
	Handle(ctx context.Context, delivery amqp.Delivery) error
	GetQueueName() string
}

type BookingNotification struct {
	UserEmail string `json:"user_email"`
	BookingID string `json:"booking_id"`
	Amount    string `json:"amount"`
	Timestamp string `json:"timestamp,omitempty"`
}

type EmailHandler struct {
	queueName   string
	emailSender email.EmailSender
	maxRetries  int
	stats       *HandlerStats
	metrics     *metrics.Metrics
}

type HandlerStats struct {
	TotalReceived   int64
	TotalSuccess    int64
	TotalFailed     int64
	TotalRetried    int64
	LastProcessedAt time.Time
}

func NewEmailHandler(queueName string, emailSender email.EmailSender, maxRetries int) *EmailHandler {
	return &EmailHandler{
		queueName:   queueName,
		emailSender: emailSender,
		maxRetries:  maxRetries,
		stats:       &HandlerStats{},
		metrics:     metrics.GetMetrics(),
	}
}

func (h *EmailHandler) Handle(ctx context.Context, delivery amqp.Delivery) error {
	// Update stats
	h.stats.TotalReceived++
	h.stats.LastProcessedAt = time.Now()

	// Record metric
	h.metrics.RabbitMQMessagesReceived.Inc()

	log.Printf("📨 [%s] Processing message: %s", h.queueName, string(delivery.Body))

	// Parse message
	var notification BookingNotification
	if err := json.Unmarshal(delivery.Body, &notification); err != nil {
		log.Printf("❌ Failed to parse JSON: %v", err)
		h.stats.TotalFailed++
		h.metrics.RabbitMQConsumeErrors.Inc()
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Validate
	if err := h.validateNotification(notification); err != nil {
		log.Printf("❌ Validation failed: %v", err)
		h.stats.TotalFailed++
		h.metrics.EmailFailedTotal.Inc()
		return err
	}

	// Generate email content
	subject := h.generateSubject(notification)
	body := h.generateBody(notification)

	// Send email with retry with metrics
	start := time.Now()
	err := h.emailSender.SendWithRetry(
		notification.UserEmail,
		subject,
		body,
		false,
		h.maxRetries,
	)
	duration := time.Since(start).Seconds()

	if err != nil {
		log.Printf("❌ Failed to send email for booking %s: %v", notification.BookingID, err)
		h.stats.TotalFailed++
		h.metrics.EmailFailedTotal.Inc()
		h.metrics.RabbitMQConsumeErrors.Inc()
		return fmt.Errorf("email sending failed: %w", err)
	}

	// Record success metrics
	h.metrics.EmailSentTotal.Inc()
	h.metrics.RabbitMQMessagesSent.Inc()

	// Record duration (you might want to add a histogram for email sending)
	log.Printf("✅ Email sent in %.2f seconds", duration)

	log.Printf("✅ [%s] Email sent successfully for booking: %s to: %s",
		h.queueName, notification.BookingID, notification.UserEmail)

	h.stats.TotalSuccess++
	return nil
}

func (h *EmailHandler) GetQueueName() string {
	return h.queueName
}

func (h *EmailHandler) GetStats() HandlerStats {
	return *h.stats
}

func (h *EmailHandler) validateNotification(notif BookingNotification) error {
	if notif.UserEmail == "" {
		return fmt.Errorf("user_email is required")
	}

	if !h.emailSender.ValidateEmail(notif.UserEmail) {
		return fmt.Errorf("invalid email format: %s", notif.UserEmail)
	}

	if notif.BookingID == "" {
		return fmt.Errorf("booking_id is required")
	}

	if notif.Amount == "" {
		return fmt.Errorf("amount is required")
	}

	return nil
}

func (h *EmailHandler) generateSubject(notif BookingNotification) string {
	return fmt.Sprintf("Booking Confirmation - %s", notif.BookingID)
}

func (h *EmailHandler) generateBody(notif BookingNotification) string {
	return fmt.Sprintf(`Dear Customer,

Your parking booking has been successfully confirmed!

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Booking Details:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✓ Booking ID: %s
✓ Total Amount: %s Tenge
✓ Status: Confirmed
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

You can view your active bookings in the mobile app.

Thank you for choosing Smart Parking!

Best regards,
Smart Parking Support Team
───────────────────────────────────
This is an automated message, please do not reply.
`, notif.BookingID, notif.Amount)
}
