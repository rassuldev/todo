package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type EmailSender struct {
	smtpHost string
	smtpPort string
	username string
	password string
	from     string
}

func NewEmailSender() *EmailSender {
	return &EmailSender{
		smtpHost: getEnv("SMTP_HOST", "smtp.gmail.com"),
		smtpPort: getEnv("SMTP_PORT", "587"),
		username: getEnv("SMTP_USERNAME", ""),
		password: getEnv("SMTP_PASSWORD", ""),
		from:     getEnv("SMTP_FROM", "noreply@todo.com"),
	}
}

func (s *EmailSender) SendEmail(to, subject, body string) error {
	// For demonstration purposes, we'll just log the email
	// In production, you would use a proper SMTP server or email service
	log.Printf("Sending email to: %s, subject: %s, body: %s", to, subject, body)

	if s.username == "" || s.password == "" {
		log.Println("SMTP credentials not configured, email not sent (simulation mode)")
		return nil
	}

	// Prepare email message
	message := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", s.from, to, subject, body))

	// Authentication
	auth := smtp.PlainAuth("", s.username, s.password, s.smtpHost)

	// Send email
	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	err := smtp.SendMail(addr, auth, s.from, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
