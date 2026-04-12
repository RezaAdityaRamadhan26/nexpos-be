package http

import (
	"net/http"
	"nexpos-be/internal/delivery/http/middleware"
	"nexpos-be/internal/models"
	"nexpos-be/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	usecase *usecase.TransactionUsecase
}

func NewTransactionHandler(r *gin.RouterGroup, usecase *usecase.TransactionUsecase) {
	handler := &TransactionHandler{usecase: usecase}

	api := r.Group("/transactions")
	{
		api.POST("/", middleware.AuthMiddleware(), handler.Checkout)
		api.GET("/", middleware.AuthMiddleware(), handler.GetAll)
		api.GET("/:id", middleware.AuthMiddleware(), handler.GetByID)
		api.GET("/dashboard", middleware.AuthMiddleware(), handler.GetDashboard)
	}
}

func (h *TransactionHandler) getIDsFromContext(c *gin.Context) (uint, uint, error) {
	userIDStr, existsUser := c.Get("user.id")
	if !existsUser {
		return 0, 0, http.ErrNoCookie
	}

	storeIDStr, existsStore := c.Get("store.id")
	if !existsStore {
		return 0, 0, http.ErrNoCookie
	}

	var UserID, StoreID uint
	if idFloat, ok := userIDStr.(float64); ok {
		UserID = uint(idFloat)
	} else {
		return 0, 0, http.ErrNotSupported
	}

	if idFloat, ok := storeIDStr.(float64); ok {
		StoreID = uint(idFloat)
	} else {
		return 0, 0, http.ErrNotSupported
	}

	return UserID, StoreID, nil
}

func (h *TransactionHandler) Checkout(c *gin.Context) {
	var trx models.Transaction

	if err := c.ShouldBindJSON(&trx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format data salah " + err.Error()})
		return
	}

	userID, storeID, err := h.getIDsFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tidak terautentikasi"})
		return
	}

	trx.UserID = userID
	trx.StoreID = storeID

	if err := h.usecase.Checkout(&trx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal transaksi: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "transaksi berhasil diproses!",
		"data":    trx,
	})
}

func (h *TransactionHandler) GetAll(c *gin.Context) {
	_, storeID, err := h.getIDsFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tidak terautentikasi"})
		return
	}

	transactions, err := h.usecase.GetTransactionHistory(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil riwayat transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "berhasil mengambil riwayat transaksi",
		"data":    transactions,
	})
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id transaksi tidak valid"})
		return
	}

	_, storeID, errCtx := h.getIDsFromContext(c)
	if errCtx != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tidak terautentikasi"})
		return
	}

	transaction, err := h.usecase.GetTransactionReceipt(id, storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaksi tidak ditemukan atau bukan milik toko ini"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "berhasil mengambil data transaksi",
		"data":    transaction,
	})
}

func (h *TransactionHandler) GetDashboard(c *gin.Context) {
	_, storeID, err := h.getIDsFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tidak terautentikasi"})
		return
	}

	dashboardData, err := h.usecase.GetDashboardData(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {"error": "gagal mengambil data dashboard" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "berhasil mengambil data dashboard",
		"data": dashboardData,
	})
	
}