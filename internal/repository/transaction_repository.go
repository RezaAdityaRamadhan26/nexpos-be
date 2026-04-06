package repository

import (
	"errors"
	"nexpos-be/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(trx *models.Transaction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(trx).Error; err != nil {
			return err
		}

		for _, detail := range trx.Details {
			var product models.Product

			if err := tx.First(&product, detail.ProductID).Error; err != nil {
				return errors.New("produk dengan id tersebut tidak ditemukan")
			}

			if product.Stock < detail.Quantity {
				return errors.New("stok '" + product.Name + "' tidak mencukupi")
			}

			product.Stock -= detail.Quantity

			if err := tx.Save(&product).Error; err != nil {
				return err
			}
		}

		return nil
	})
}