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
		var grandTotal float64

		for i, detail := range trx.Details {
			var product models.Product

			if err := tx.Where("id = ? AND store_id = ?", detail.ProductID, trx.StoreID).First(&product).Error; err != nil {
				return errors.New("produk tidak ditemukan atau bukan milik toko ini")
			}

			if product.Stock < detail.Quantity {
				return errors.New("stok '" + product.Name + "' tidak mencukupi")
			}

			trx.Details[i].Price = product.Price 
			subTotal := product.Price * float64(detail.Quantity)
			trx.Details[i].SubTotal = subTotal
			
			grandTotal += subTotal

			product.Stock -= detail.Quantity

			if err := tx.Save(&product).Error; err != nil {
				return err
			}
		}

		trx.TotalAmount = grandTotal

		if err := tx.Create(trx).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *TransactionRepository) FindAllByStoreID(storeID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := r.db.Where("store_id = ?", storeID).
			Preload("Details").
			Order("created_at desc").
			Find(&transactions).Error
			
			return transactions, err
}

func (r *TransactionRepository) FindByIDAndStoreID(id int, storeID uint) (*models.Transaction, error) {
	var transaction models.Transaction

	err := r.db.Where("id = ? AND store_id = ?", id, storeID).
			Preload("Details").
			Preload("Details.Product").
			Preload("User").
			Preload("Store").
			First(&transaction).Error
			
			return &transaction, err
}

func (r *TransactionRepository) GetDashboardStats(storeID uint) (float64, int64, error) {
	var totalRevenue float64
	var totalTransactions int64

	errCount := r.db.Model(&models.Transaction{}).Where("store_id = ?", storeID).Count(&totalTransactions).Error
	if errCount != nil {
		return 0, 0, errCount
	}

	errSum := r.db.Model(&models.Transaction{}).
	Where("store_id = ?", storeID).
	Select("COALESCE(SUM(total_amount), 0)").
	Scan(&totalRevenue).Error

	if errSum != nil {
		return 0, 0, errSum
	}

	return totalRevenue, totalTransactions, nil
}