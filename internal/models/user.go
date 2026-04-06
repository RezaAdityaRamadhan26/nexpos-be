package models

import "time"

type User struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(100);not null" json:"name"`
	Email string `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"` // json:"-" password ga ikut kekirim ke fe nextjs
	Role string `gorm:"type:varchar(20);default:'kasir'" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}