package usecase

import (
	"errors"
	"nexpos-be/internal/models"
	"nexpos-be/internal/repository"
)

type TransactionUsecase struct {
	repo *repository.TransactionRepository
}

type DashboardResponse struct {
	TotalRevenue 		float64 		`json:"total_revenue"`
	TotalTransactions 	int64 			`json:"total_transactions"`
}

func NewTransactionUsecase(repo *repository.TransactionRepository) *TransactionUsecase {
	return &TransactionUsecase{repo: repo}
}

func (u *TransactionUsecase) Checkout(trx *models.Transaction) error {
	if len(trx.Details) == 0 {
		return errors.New("keranjang tidak boleh kosong!")
	}
	trx.Status = "completed"

	return u.repo.CreateTransaction(trx)
}

func (u *TransactionUsecase) GetTransactionHistory(storeID uint) ([]models.Transaction, error) {
	return u.repo.FindAllByStoreID(storeID)
}

func (u *TransactionUsecase) GetTransactionReceipt(id int, storeID uint) (*models.Transaction, error) {
	return u.repo.FindByIDAndStoreID(id, storeID)
}

func (u *TransactionUsecase) GetDashboardData(storeID uint) (DashboardResponse, error) {
	revenue, count, err := u.repo.GetDashboardStats(storeID)
	if err != nil{
		return DashboardResponse{}, err
	}

	return DashboardResponse{
		TotalRevenue: revenue,
		TotalTransactions: count,
	}, nil
}