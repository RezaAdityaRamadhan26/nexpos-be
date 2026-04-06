package usecase

import (
	"errors"
	"nexpos-be/internal/models"
	"nexpos-be/internal/repository"
)

type TransactionUsecase struct {
	repo *repository.TransactionRepository
}

func NewTransactionUsecase(repo *repository.TransactionRepository) *TransactionUsecase {
	return &TransactionUsecase{repo: repo}
}

func (u *TransactionUsecase) Checkout(trx *models.Transaction) error {
	if len(trx.Details) == 0 {
		return errors.New("keranjang tidak boleh kosong!")
	}

	trx.Status = "completed"

	var TotalAmount float64 
	for i, detail := range trx.Details {
		subTotal := detail.Price * float64(detail.Quantity) 

		trx.Details[i].SubTotal = subTotal
	}

	trx.TotalAmount = TotalAmount

	return u.repo.CreateTransaction(trx)
}