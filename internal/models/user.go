package models

import "time"

type User struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	StoreID uint `gorm:"not null" json:"store_id"`
	Name string `gorm:"type:varchar(100);not null" json:"name"`
	Email string `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"` // json:"-" password ga ikut kekirim ke fe nextjs
	Role string `gorm:"type:varchar(20);default:'kasir'" json:"role"`
	IsVerified    bool      `gorm:"default:false" json:"is_verified"`
	OTPCode       string    `gorm:"type:varchar(6)" json:"-"` // Jangan kirim OTP ke JSON Frontend
	OTPExpiration time.Time `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}