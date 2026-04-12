	package usecase

	import (
		"nexpos-be/internal/models"
		"nexpos-be/internal/repository"

		"errors"
	)

	type ProductUsecase struct {
		repo *repository.ProductRepository
	}

	func NewProductUsecase (repo *repository.ProductRepository) *ProductUsecase {
		return &ProductUsecase{repo: repo}
	}

	func (u *ProductUsecase) CreateProduct(product *models.Product) error {
		return u.repo.Create(product)
	}

	func (u *ProductUsecase) GetAllProducts(storeID uint) ([]models.Product, error) {
		return u.repo.FindByStoreID(storeID)
	}

	func (u *ProductUsecase) UpdateProducts(id int, storeID uint, updatedData *models.Product) error {
		product, err := u.repo.FindByIDAndStoreID(id, storeID)
		if err != nil {
			return errors.New("produk tidak ditemukan")
		}

		product.Name = updatedData.Name
		product.SKU = updatedData.SKU
		product.Price = updatedData.Price
		product.Stock = updatedData.Stock
		product.Category = updatedData.Category
		product.Description = updatedData.Description
		product.ImageURL = updatedData.ImageURL

		return u.repo.Update(product)
	}
	
	func (u *ProductUsecase) DeleteProducts(id int, storeID uint) error {
		_, err := u.repo.FindByIDAndStoreID(id, storeID)
		if err != nil {
			return errors.New("produk tidak ditemukan")
		}

		return u.repo.Delete(id, storeID)
	}