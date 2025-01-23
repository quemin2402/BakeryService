package models

import "time"

type Role struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`
}

type User struct {
	ID                uint    `gorm:"primaryKey" json:"id"`
	Email             string  `gorm:"unique;not null" json:"email"`
	Password          string  `gorm:"not null" json:"-"`
	RoleID            uint    `json:"role_id"`
	Role              Role    `gorm:"foreignKey:RoleID" json:"role"`
	EmailVerified     bool    `gorm:"default:false" json:"email_verified"`
	ConfirmationToken *string `gorm:"unique" json:"confirmation_token"`

	FirstName string    `gorm:"default:''" json:"first_name"`
	LastName  string    `gorm:"default:''" json:"last_name"`
	BirthYear *int      `json:"birth_year,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Orders []Order `gorm:"foreignKey:UserID" json:"orders"`
}
