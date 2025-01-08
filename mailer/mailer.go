package mailer

import (
	"github.com/go-gomail/gomail"
	"log"
)

func SendEmail(smtpHost string, smtpPort int, senderEmail, senderPassword, recipient, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error sending email to %s: %v", recipient, err)
		return err
	}

	log.Printf("Email sent to %s successfully", recipient)
	return nil
}
