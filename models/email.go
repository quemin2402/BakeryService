package models

import "time"

type Email struct {
	ID             int       `json:"id"`
	RecipientEmail string    `json:"recipient_email"`
	Subject        string    `json:"subject"`
	Body           string    `json:"body"`
	SentAt         time.Time `json:"sent_at"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}
