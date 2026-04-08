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

func (u *ProductUsecase) GetAllProducts() ([] models.Product, error) {
	return u.repo.GetAll()
}

func (u *ProductUsecase) UpdateProducts(id int, updatedData	*models.Product) error {
	product, err := u.repo.FindByID(id)
	if err != nil {
		return errors.New("produk tidak ditemukan!")
	}
	product.Name = updatedData.Name
	product.Price = updatedData.Price
	product.Stock = updatedData.Stock

	return u.repo.Update(product)
}

func (u *ProductUsecase) DeleteProducts(id int) error {
	_, err := u.repo.FindByID(id)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}

	return u.repo.Delete(id)
}