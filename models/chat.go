package models

import "time"

type Chat struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint `gorm:"unique"`
	IsActive bool `gorm:"default:true"`
}

type Message struct {
	ID        uint `gorm:"primaryKey"`
	ChatID    uint `gorm:"index"`
	UserID    uint
	Content   string
	Timestamp time.Time `gorm:"autoCreateTime"`
}
