package handlers

import (
	"BakeryService/config"
	_ "BakeryService/logger"
	"BakeryService/mailer"
	"encoding/json"
	_ "encoding/json"
	_ "log"
	"net/http"
)

type EmailRequest struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func SendEmailHandler(config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		recipient := r.FormValue("recipient")
		subject := r.FormValue("subject")
		body := r.FormValue("body")

		if recipient == "" || subject == "" || body == "" {
			http.Error(w, "Recipient, Subject, and Body are required", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("attachment")
		if err != nil && err != http.ErrMissingFile {
			http.Error(w, "Failed to read attachment", http.StatusInternalServerError)
			return
		}
		defer func() {
			if file != nil {
				file.Close()
			}
		}()

		err = mailer.SendEmailWithAttachment(
			config.SMTPHost,
			config.SMTPPort,
			config.SenderEmail,
			config.SenderPassword,
			recipient,
			subject,
			body,
			file,
			header.Filename,
		)
		if err != nil {
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Email sent successfully"})
	}
}
