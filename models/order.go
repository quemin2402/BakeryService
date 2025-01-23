package models

import "time"

type Order struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `json:"user_id"`
	User       User           `gorm:"foreignKey:UserID" json:"user"`
	OrderDate  time.Time      `gorm:"autoCreateTime" json:"order_date"`
	Products   []OrderProduct `gorm:"foreignKey:OrderID" json:"products"`
	TotalPrice float64        `json:"total_price"`
}
