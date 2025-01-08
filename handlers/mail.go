package handlers

import (
	"BakeryService/config"
	"BakeryService/logger"
	"BakeryService/mailer"
	"encoding/json"
	"net/http"
)

type EmailRequest struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func SendEmailHandler(config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var emailReq EmailRequest
		if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
			logger.LogEntry(map[string]interface{}{
				"action": "decode_email_request",
				"error":  err.Error(),
			})
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if emailReq.Recipient == "" || emailReq.Subject == "" || emailReq.Body == "" {
			logger.LogEntry(map[string]interface{}{
				"action": "validate_email_request",
				"error":  "Missing required fields",
			})
			http.Error(w, "Recipient, Subject, and Body are required", http.StatusBadRequest)
			return
		}

		err := mailer.SendEmail(
			config.SMTPHost,
			config.SMTPPort,
			config.SenderEmail,
			config.SenderPassword,
			emailReq.Recipient,
			emailReq.Subject,
			emailReq.Body,
		)
		if err != nil {
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Email sent successfully"})
	}
}
