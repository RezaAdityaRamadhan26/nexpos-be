package models

import "time"

type Store struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"type:varchar(100);not null" json:"name"`
	Address          string    `gorm:"type:text" json:"address"`
	SubscriptionTier string    `gorm:"type:varchar(20);default:'starter'" json:"subscription_tier"` // starter(gratis), premium(menengah), enterprise(mahal)
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// relasi, 1 toko punya banyak kasir, produk sama transaksi
	Users        []User        `gorm:"foreignKey:StoreID" json:"-"`
	Products     []Product     `gorm:"foreignKey:StoreID" json:"-"`
	Transactions []Transaction `gorm:"foreignKey:StoreID" json:"-"`
}