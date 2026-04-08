package http 

import (
	"net/http"
	"nexpos-be/internal/models"
	"nexpos-be/internal/usecase"
	
	"github.com/gin-gonic/gin"
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

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func (h *UserHandler) RegisterStore(c * gin.Context) {
	var req usecase.RegisterStoreRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.usecase.RegisterStore(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Toko Berhasil Didaftarkan, silahkan cek email untuk kode verifikasi",
	})
}

func (h *UserHandler) VerifyOTP(c *gin.Context) {
	var req usecase.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.VerifyOTP(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verifikasi barhasil!, akun dan toko sekarang aktif",
	})
}

func NewUserHandler (r *gin.Engine, usecase *usecase.UserUsecase) {
	handler := &UserHandler{usecase: usecase}

	api := r.Group("/api/users")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
		api.POST("/register-store", handler.RegisterStore)
		api.POST("/verify-otp", handler.VerifyOTP)
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var registerData struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Name:     registerData.Name,
		Email:    registerData.Email,
		Password: registerData.Password,
		Role:     registerData.Role,
	}

	if err := h.usecase.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mendaftar " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registrasi berhasil!"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var LoginData struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&LoginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.usecase.Login(LoginData.Email, LoginData.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Berhasil!",
		"token": token,
	})
}