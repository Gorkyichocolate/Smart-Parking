// internal/email/sender.go
package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"
)

type EmailSender interface {
	Send(to, subject, body string, isHTML bool) error
	SendWithRetry(to, subject, body string, isHTML bool, maxRetries int) error
	ValidateEmail(email string) bool
}

type SMTPSender struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewSMTPSender(host string, port int, username, password, from string) *SMTPSender {
	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SMTPSender) Send(to, subject, body string, isHTML bool) error {
	log.Printf("📧 Preparing email to: %s, subject: %s", to, subject)

	if !s.ValidateEmail(to) {
		return fmt.Errorf("invalid email address: %s", to)
	}

	contentType := "text/plain"
	if isHTML {
		contentType = "text/html"
	}

	headers := make(map[string]string)
	headers["From"] = s.from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("%s; charset=utf-8", contentType)

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	if err := s.sendViaSMTP(addr, auth, to, message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("✅ Email sent successfully to: %s", to)
	return nil
}

func (s *SMTPSender) SendWithRetry(to, subject, body string, isHTML bool, maxRetries int) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := s.getRetryDelay(attempt)
			log.Printf("🔄 Retry %d/%d after %v", attempt, maxRetries, delay)
			time.Sleep(delay)
		}

		if err := s.Send(to, subject, body, isHTML); err != nil {
			lastErr = err
			log.Printf("❌ Attempt %d failed: %v", attempt+1, err)

			if !s.isTemporaryError(err) {
				return err
			}
			continue
		}

		if attempt > 0 {
			log.Printf("✅ Email sent on retry %d", attempt)
		}
		return nil
	}

	return fmt.Errorf("failed after %d attempts: %w", maxRetries+1, lastErr)
}

func (s *SMTPSender) sendViaSMTP(addr string, auth smtp.Auth, to, message string) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()

	if err = client.Hello(s.host); err != nil {
		return fmt.Errorf("HELO failed: %w", err)
	}

	if err = client.StartTLS(&tls.Config{
		ServerName: s.host,
	}); err != nil {
		return fmt.Errorf("STARTTLS failed: %w", err)
	}

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("auth failed: %w", err)
	}

	if err = client.Mail(s.from); err != nil {
		return fmt.Errorf("sender failed: %w", err)
	}

	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("recipient failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("data failed: %w", err)
	}
	defer w.Close()

	if _, err = w.Write([]byte(message)); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	return nil
}

func (s *SMTPSender) ValidateEmail(email string) bool {
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

func (s *SMTPSender) getRetryDelay(attempt int) time.Duration {
	delay := time.Duration(1<<uint(attempt)) * time.Second
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}
	return delay
}

func (s *SMTPSender) isTemporaryError(err error) bool {
	errStr := strings.ToLower(err.Error())
	temporaryErrors := []string{
		"connection refused",
		"timeout",
		"temporary failure",
		"try again",
		"too many connections",
		"rate limit",
	}
	for _, tempErr := range temporaryErrors {
		if strings.Contains(errStr, tempErr) {
			return true
		}
	}
	return false
}
