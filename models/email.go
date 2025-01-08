package models

import "time"

type Email struct {
	ID             uint   `gorm:"primaryKey"`
	RecipientEmail string `gorm:"not null"`
	Subject        string `gorm:"not null"`
	Body           string `gorm:"not null"`
	SentAt         *time.Time
	Status         string    `gorm:"default:pending"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}
