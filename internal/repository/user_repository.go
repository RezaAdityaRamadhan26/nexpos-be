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

func (r *UserRepository) CreateStorAndOwner(store *models.Store, owner *models.User) error {
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