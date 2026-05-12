package usecase

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/GorkyiChocolate/smart-parking/services/notification-service/internal/domain"
)

type NotifierUseCase struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	smtpFrom     string
	maxRetries   int
}

func NewNotifierUseCase(
	host string,
	port int,
	username, password, from string,
	maxRetries int,
) *NotifierUseCase {
	return &NotifierUseCase{
		smtpHost:     host,
		smtpPort:     port,
		smtpUsername: username,
		smtpPassword: password,
		smtpFrom:     from,
		maxRetries:   maxRetries,
	}
}

func (uc *NotifierUseCase) SendEmail(notif domain.EmailNotification) error {
	log.Printf("📧 Sending email to: %s, subject: %s", notif.To, notif.Subject)

	if !uc.ValidateEmail(notif.To) {
		return fmt.Errorf("invalid email address: %s", notif.To)
	}

	contentType := "text/plain"
	if notif.IsHTML {
		contentType = "text/html"
	}

	headers := make(map[string]string)
	headers["From"] = uc.smtpFrom
	headers["To"] = notif.To
	headers["Subject"] = notif.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("%s; charset=utf-8", contentType)

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + notif.Body

	auth := smtp.PlainAuth("", uc.smtpUsername, uc.smtpPassword, uc.smtpHost)
	addr := fmt.Sprintf("%s:%d", uc.smtpHost, uc.smtpPort)

	var err error
	for attempt := 0; attempt <= uc.maxRetries; attempt++ {
		if attempt > 0 {
			delay := uc.GetRetryDelay(attempt)
			log.Printf("🔄 Retry %d/%d after %v", attempt, uc.maxRetries, delay)
			time.Sleep(delay)
		}

		// Передаем email получателя как параметр
		err = uc.sendViaSMTP(addr, auth, notif.To, message)
		if err == nil {
			log.Printf("✅ Email sent successfully to: %s", notif.To)
			return nil
		}

		log.Printf("❌ Attempt %d failed: %v", attempt+1, err)
	}

	return fmt.Errorf("failed to send email after %d attempts: %w", uc.maxRetries+1, err)
}

// Исправленный метод - добавляем параметр to
func (uc *NotifierUseCase) sendViaSMTP(addr string, auth smtp.Auth, to, message string) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to dial SMTP: %w", err)
	}
	defer client.Close()

	// Отправляем HELO/EHLO
	if err = client.Hello(uc.smtpHost); err != nil {
		return fmt.Errorf("failed to say HELO: %w", err)
	}

	// STARTTLS
	if err = client.StartTLS(&tls.Config{
		ServerName: uc.smtpHost,
		MinVersion: tls.VersionTLS12,
	}); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// Аутентификация
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	// Отправитель
	if err = client.Mail(uc.smtpFrom); err != nil {
		return fmt.Errorf("failed to set from: %w", err)
	}

	// Получатель - используем параметр to
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Тело письма
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer w.Close()

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (uc *NotifierUseCase) ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) < 3 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func (uc *NotifierUseCase) GetRetryDelay(attempt int) time.Duration {
	// Exponential backoff: 1s, 2s, 4s, 8s...
	delay := time.Duration(1<<uint(attempt)) * time.Second
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}
	return delay
}
