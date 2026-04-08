package models

import "time"

type Transaction struct {
	ID            uint                `gorm:"primaryKey" json:"id"`
	StoreID       uint                `gorm:"not null" json:"store_id"`
	UserID        uint                `gorm:"not null" json:"user_id"` // id kasir
	TotalAmount   float64             `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	PaymentMethod string              `gorm:"type:varchar(50)" json:"payment_method"` // contoh : tunai, qris, transfer
	Status        string              `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, paid, completed, cancelled
	Details       []TransactionDetail `gorm:"foreignKey:TransactionID" json:"details"` // relasi ke tabel detail
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}

type TransactionDetail struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TransactionID uint      `gorm:"not null" json:"transaction_id"`
	ProductID     uint      `gorm:"not null" json:"product_id"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	Price         float64   `gorm:"type:decimal(12,2);not null" json:"price"` // harga satuan pas transaksi
	SubTotal      float64   `gorm:"type:decimal(12,2);not null" json:"sub_total"` // quantity * price
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}