	package repository

	import (
		"nexpos-be/internal/models"

		"gorm.io/gorm"
	)

	type ProductRepository struct {
		db *gorm.DB
	}

	//fungsi buat bikin "tukang" repository baru
	func NewProductRepository(db *gorm.DB) *ProductRepository {
		return &ProductRepository{db: db}
	}

	// nyimpen produk ke db
	func (r* ProductRepository) Create(product *models.Product) error {
		return r.db.Create(product).Error
	}

	// mengambil semua produk dari database
	func (r *ProductRepository) FindByStoreID(storeID uint) ([]models.Product, error) {
		var product []models.Product

		err := r.db.Where("store_id = ?", storeID).Find(&product).Error
		return product, err
	}

	func (r *ProductRepository) FindByIDAndStoreID(id int, storeID uint) (*models.Product, error) {
		var product models.Product

		err := r.db.Where("id = ? AND store_id = ?", id, storeID).First(&product).Error
		return &product, err
	}

	func (r* ProductRepository) Update(product *models.Product)  error {
		return r.db.Save(product).Error
	}

	func (r* ProductRepository) Delete(id int, storeID uint) error {
		return r.db.Where("id = ? AND store_id = ?", id, storeID).Delete(&models.Product{}).Error
	} 





