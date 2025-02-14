package models

import "time"

type Chat struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint
	IsActive bool `gorm:"default:true"`
}

type Message struct {
	ID        uint `gorm:"primaryKey"`
	ChatID    uint `gorm:"index"`
	UserID    uint
	Content   string
	Timestamp time.Time `gorm:"autoCreateTime"`
}

type MessageResponse struct {
	ID        uint
	ChatID    uint
	UserID    uint
	Content   string
	Timestamp time.Time
	Role      uint
}
