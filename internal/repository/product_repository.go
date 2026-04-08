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
func (r *ProductRepository) GetAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r* ProductRepository) FindByID(id int) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r* ProductRepository) Update(product *models.Product)  error {
	return r.db.Save(product).Error
}

func (r* ProductRepository) Delete(id int) error {
	return r.db.Delete(&models.Product{}, id).Error
} 





