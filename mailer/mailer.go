package mailer

import (
	"BakeryService/config"
	"fmt"
	"github.com/go-gomail/gomail"
	"io"
	"os"
	"path/filepath"
)

func SendEmailWithAttachment(smtpHost string, smtpPort int, senderEmail, senderPassword, recipient, subject, body string, attachment io.Reader, filename string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if attachment != nil {
		tempDir := os.TempDir()
		tempFilePath := filepath.Join(tempDir, filename)

		file, err := os.Create(tempFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, attachment)
		if err != nil {
			return err
		}

		m.Attach(tempFilePath)
		defer os.Remove(tempFilePath)
	}

	d := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendEmailWithConfirmation(cfg *config.Config, recipientEmail, confirmationToken string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.SenderEmail)
	m.SetHeader("To", recipientEmail)
	m.SetHeader("Subject", "Confirm Your Email")
	body := fmt.Sprintf(`
        <p>Thank you for registering!</p>
        <p>Please click the link below to verify your email address:</p>
        <a href="http://localhost:8080/api/auth/confirm-email?token=%s">Verify Email</a>
    `, confirmationToken)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SenderEmail, cfg.SenderPassword)

	if err := dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
