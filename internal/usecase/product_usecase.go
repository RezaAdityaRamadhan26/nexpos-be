package usecase

import (
	"nexpos-be/internal/models"
	"nexpos-be/internal/repository"
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