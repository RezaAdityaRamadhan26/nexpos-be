package http

import (
	"net/http"
	"nexpos-be/internal/delivery/http/middleware"
	"nexpos-be/internal/models"
	"nexpos-be/internal/usecase"
	"strconv"

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

func (h *UserHandler) RegisterStore(c *gin.Context) {
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

func NewUserHandler(r *gin.Engine, usecase *usecase.UserUsecase) {
	handler := &UserHandler{usecase: usecase}

	api := r.Group("/api/users")
	{	// Main Route
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
		api.POST("/register-store", handler.RegisterStore)
		api.POST("/verify-otp", handler.VerifyOTP)

		// Staff Route
		api.POST("/staff", middleware.AuthMiddleware(), middleware.OwnerOnly(), handler.RegisterStaff)
		api.GET("/staff", middleware.AuthMiddleware(), middleware.OwnerOnly(), handler.GetStaffList)
		api.PUT("/staff/:id", middleware.AuthMiddleware(), middleware.OwnerOnly(), handler.UpdateStaff)
		api.DELETE("/staff/:id", middleware.AuthMiddleware(), middleware.OwnerOnly(), handler.DeleteStaff)

		// Profile Router
		api.GET("/profile", middleware.AuthMiddleware(), handler.GetProfile)
		api.PUT("/profile", middleware.AuthMiddleware(), handler.UpdateProfile)
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
		Email    string `json:"email"`
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
		"token":   token,
	})
}

func (h *UserHandler) RegisterStaff(c *gin.Context) {
	var req usecase.RegisterStaffRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storeIDStr, _ := c.Get("store.id")
	storeID := uint(storeIDStr.(float64))

	if err := h.usecase.RegisterStaff(req, storeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mendaftarkan kasir" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "kasir berhasil ditambahkan"})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user.id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tidak terautentikasi"})
		return
	}

	var UserID uint
	switch v := userIDVal.(type) {
	case float64:
		UserID = uint(v)

	case uint:
		UserID = v
	}

	user, err := h.usecase.GetProfile(UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User Tidak Ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil Mengambil Profile",
		"data": gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"role":     user.Role,
			"store_id": user.StoreID,
		},
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDVal, _ := c.Get("user.id")
	var UserID uint
	switch v := userIDVal.(type) {
	case float64:
		UserID = uint(v)
	case uint:
		UserID = v
	}

	var req usecase.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format Data Tidak Valid"})
		return
	}

	err := h.usecase.UpdateProfile(UserID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Memperbarui Profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil Memperbarui Profile"})
}

func (h *UserHandler) GetStaffList(c *gin.Context) {
	storeID, errCtx := getStoreIDFromContext(c)
	if errCtx != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tidak terautentikasi"})
		return
	}
	staff, err := h.usecase.GetStaffList(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data Pegawai"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": staff})
}

func (h *UserHandler) UpdateStaff(c *gin.Context) {
	storeID, errCtx := getStoreIDFromContext(c)
	if errCtx != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tidak terautentikasi"})
		return
	}
	staffIDParam := c.Param("id")
	staffID, _ := strconv.Atoi(staffIDParam)

	var req usecase.UpdateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	if err := h.usecase.UpdateStaff(uint(staffID), storeID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data pegawai berhasil diperbarui"})
}

func (h *UserHandler) DeleteStaff(c *gin.Context) {
	storeID, errCtx := getStoreIDFromContext(c)
	if errCtx != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tidak terautentikasi"})
		return
	}
	staffIDParam := c.Param("id")
	staffID, _ := strconv.Atoi(staffIDParam)

	if err := h.usecase.Deletestaff(uint(staffID), storeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pegawai berhasil dihapus"})
}
