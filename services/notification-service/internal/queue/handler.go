// internal/queue/handler.go
package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"notification-service/internal/email"
)

// MessageHandler интерфейс обработчика сообщений
type MessageHandler interface {
	Handle(ctx context.Context, delivery amqp.Delivery) error
	GetQueueName() string
}

// BookingNotification структура сообщения из Payment Service
type BookingNotification struct {
	UserEmail string `json:"user_email"`
	BookingID string `json:"booking_id"`
	Amount    string `json:"amount"`
	Timestamp string `json:"timestamp,omitempty"`
}

// EmailHandler обработчик для email уведомлений
type EmailHandler struct {
	queueName   string
	emailSender email.EmailSender
	maxRetries  int
	stats       *HandlerStats
}

// HandlerStats статистика обработки
type HandlerStats struct {
	TotalReceived   int64
	TotalSuccess    int64
	TotalFailed     int64
	TotalRetried    int64
	LastProcessedAt time.Time
}

// NewEmailHandler создает новый обработчик email
func NewEmailHandler(queueName string, emailSender email.EmailSender, maxRetries int) *EmailHandler {
	return &EmailHandler{
		queueName:   queueName,
		emailSender: emailSender,
		maxRetries:  maxRetries,
		stats:       &HandlerStats{},
	}
}

// Handle обрабатывает сообщение из очереди
func (h *EmailHandler) Handle(ctx context.Context, delivery amqp.Delivery) error {
	// Обновляем статистику
	h.stats.TotalReceived++
	h.stats.LastProcessedAt = time.Now()

	log.Printf("📨 [%s] Processing message: %s", h.queueName, string(delivery.Body))

	// Парсим сообщение
	var notification BookingNotification
	if err := json.Unmarshal(delivery.Body, &notification); err != nil {
		log.Printf("❌ Failed to parse JSON: %v", err)
		h.stats.TotalFailed++
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Валидация обязательных полей
	if err := h.validateNotification(notification); err != nil {
		log.Printf("❌ Validation failed: %v", err)
		h.stats.TotalFailed++
		return err
	}

	// Формируем email
	subject := h.generateSubject(notification)
	body := h.generateBody(notification)
	isHTML := false // Можно сделать HTML шаблоны

	// Отправляем email с retry
	err := h.emailSender.SendWithRetry(
		notification.UserEmail,
		subject,
		body,
		isHTML,
		h.maxRetries,
	)

	if err != nil {
		log.Printf("❌ Failed to send email for booking %s: %v", notification.BookingID, err)
		h.stats.TotalFailed++
		return fmt.Errorf("email sending failed: %w", err)
	}

	log.Printf("✅ [%s] Email sent successfully for booking: %s to: %s",
		h.queueName, notification.BookingID, notification.UserEmail)

	h.stats.TotalSuccess++
	return nil
}

// GetQueueName возвращает имя очереди
func (h *EmailHandler) GetQueueName() string {
	return h.queueName
}

// GetStats возвращает статистику обработки
func (h *EmailHandler) GetStats() HandlerStats {
	return *h.stats
}

// validateNotification валидирует входящее уведомление
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

// generateSubject генерирует тему письма
func (h *EmailHandler) generateSubject(notif BookingNotification) string {
	return fmt.Sprintf("Booking Confirmation - %s", notif.BookingID)
}

// generateBody генерирует тело письма
func (h *EmailHandler) generateBody(notif BookingNotification) string {
	body := fmt.Sprintf(`Dear Customer,

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

	// Добавляем timestamp если есть
	if notif.Timestamp != "" {
		body += fmt.Sprintf("\nProcessed at: %s", notif.Timestamp)
	}

	return body
}

// GetHTMLBody генерирует HTML версию письма (для будущих улучшений)
func (h *EmailHandler) GetHTMLBody(notif BookingNotification) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .details { background: white; padding: 15px; margin: 10px 0; border-left: 4px solid #4CAF50; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #777; }
        .status { color: #4CAF50; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Smart Parking</h2>
            <h3>Booking Confirmation</h3>
        </div>
        <div class="content">
            <p>Dear Customer,</p>
            <p>Your parking booking has been successfully confirmed!</p>
            <div class="details">
                <p><strong>Booking ID:</strong> %s</p>
                <p><strong>Total Amount:</strong> %s Tenge</p>
                <p><strong>Status:</strong> <span class="status">✓ Confirmed</span></p>
            </div>
            <p>You can view your active bookings in the mobile app.</p>
            <p>Thank you for choosing Smart Parking!</p>
        </div>
        <div class="footer">
            <p>This is an automated message, please do not reply.</p>
            <p>&copy; 2024 Smart Parking System</p>
        </div>
    </div>
</body>
</html>`, notif.BookingID, notif.Amount)
}
