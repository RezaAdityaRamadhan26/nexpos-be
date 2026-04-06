package models

import "time"

type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	SKU       string    `gorm:"type:varchar(50);unique;not null" json:"sku"` // Kode unik barang
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock     int       `gorm:"not null;default:0" json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}