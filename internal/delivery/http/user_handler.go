package http 

import (
	"net/http"
	"nexpos-be/internal/models"
	"nexpos-be/internal/usecase"
	
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler (r *gin.Engine, usecase *usecase.UserUsecase) {
	handler := &UserHandler{usecase: usecase}

	api := r.Group("/api/users")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
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