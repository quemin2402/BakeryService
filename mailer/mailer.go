package mailer

import (
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
