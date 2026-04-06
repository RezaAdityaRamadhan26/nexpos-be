package usecase 

import (
	"errors"
	"os"
	"time"

	"nexpos-be/internal/models"
	"nexpos-be/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo *repository.UserRepository
}

func NewUserUsecase(repo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) Register(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}  
	user.Password = string(hashedPassword)

	return u.repo.Create(user)
}

func (u *UserUsecase) Login(email, password string) (string, error) {
	user, err := u.repo.FindByEmail(email)
	if err !=  nil {
		return "", errors.New("email atau password salah!")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user.id": user.ID,
		"role": user.Role,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // berlaku 24 jam	
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}



