package repository

import (
	"nexpos-be/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) CreateStoreAndOwner(store *models.Store, owner *models.User) error {
	tx := r.db.Begin()

	if err := tx.Create(store).Error; err != nil {
		tx.Rollback()
		return err
	}

	owner.StoreID = store.ID

	if err := tx.Create(owner).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *UserRepository) FindByID(id uint) (models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *UserRepository) GetStaffByID(storeID uint) ([]models.User, error) {
	var staff []models.User
	err := r.db.Where("store_id = ? AND role = ?", storeID, "cashier").Find(&staff).Error
	return staff, err
} 

func (r *UserRepository) Delete(user *models.User) error {
	return r.db.Delete(user).Error
}