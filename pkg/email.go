package pkg

import (
	"fmt"
	"net/smtp"
)

// SendEmail sends an email using SMTP
func SendEmail(to, subject, body string) error {
	from := "shikubuinjila@gmail.com"
	password := "your-email-password"

	// SMTP server configuration.
	smtpHost := "smtp.example.com"
	smtpPort := "587"

	// Message format
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
