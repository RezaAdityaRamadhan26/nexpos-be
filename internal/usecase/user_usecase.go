package usecase

import (
	"errors"
	"log"
	"os"
	"time"

	"nexpos-be/internal/models"
	"nexpos-be/internal/repository"
	"nexpos-be/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterStoreRequest struct {
	StoreName    string `json:"store_name" binding:"required"`
	StoreAddress string `json:"store_address"`
	OwnerName    string `json:"owner_name" binding:"required"`
	OwnerEmail   string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
}

type VerifyOTPRequest struct {
	Email   string `json:"email" binding:"required,email"`
	OTPCode string `json:"otp_code" binding:"required"`	
}

type UserUsecase struct {
	repo *repository.UserRepository
}

func (u *UserUsecase) RegisterStore(req RegisterStoreRequest) error {
	_, err := u.repo.FindByEmail(req.OwnerEmail)
	if err == nil {
		return errors.New("email sudah digunakan")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	store := models.Store{
		Name: req.StoreName,
		Address: req.StoreAddress,
		SubscriptionTier: "starter",
	}

	otpCode := utils.GenerateOTP()

	owner := models.User{
		Name: req.OwnerName,
		Email: req.OwnerEmail,
		Password: string(hashedPassword),
		Role: "owner",
		IsVerified: false,
		OTPCode: otpCode,
		OTPExpiration: time.Now().Add(10 * time.Minute),
	}

	if err := u.repo.CreateStorAndOwner(&store, &owner); err != nil {
		return err
	}

go func(email string, otp string) {
		errMail := utils.SendOTPEmailAPI(email, otp)
		if errMail != nil {
			log.Printf("Gagal mengirim email OTP ke %s: %v\n", email, errMail)
		}
	} (owner.Email, otpCode)

	return nil
}

func (u *UserUsecase) VerifyOTP(req VerifyOTPRequest) error {
	user, err := u.repo.FindByEmail(req.Email)
	if err != nil {
		return errors.New("pengguna tidak ditemukan")
	}

	if user.IsVerified {
		return errors.New("akun sudah terverifikasi")
	}

	if time.Now().After(user.OTPExpiration) {
		return errors.New("kode OTP sudah kadaluarsa, silahkan minta ulang")
	}

	if user.OTPCode != req.OTPCode {
		return errors.New("kode OTP salah")
	}

	user.IsVerified = true
	user.OTPCode = ""

	return u.repo.Update(user)
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



