package usecase

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"notification-service/internal/domain"
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

		err = uc.sendViaSMTP(addr, auth, message)
		if err == nil {
			log.Printf("Email sent successfully to: %s", notif.To)
			return nil
		}

		log.Printf("Attempt %d failed: %v", attempt+1, err)
	}

	return fmt.Errorf("failed to send email after %d attempts: %w", uc.maxRetries+1, err)
}

func (uc *NotifierUseCase) sendViaSMTP(addr string, auth smtp.Auth, message string) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to dial SMTP: %w", err)
	}
	defer client.Close()

	if err = client.StartTLS(&tls.Config{
		ServerName: uc.smtpHost,
	}); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	if err = client.Mail(uc.smtpFrom); err != nil {
		return fmt.Errorf("failed to set from: %w", err)
	}

	if err = client.Rcpt(uc.smtpTo); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

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
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func (uc *NotifierUseCase) GetRetryDelay(attempt int) time.Duration {
	// Exponential backoff: 1s, 2s, 4s, 8s...
	return time.Duration(1<<uint(attempt)) * time.Second
}
