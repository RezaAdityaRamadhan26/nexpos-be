package http

import (
	"net/http"
	"nexpos-be/internal/models"
	"nexpos-be/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	usecase	*usecase.TransactionUsecase
}

func NewTransactionHandler (r *gin.RouterGroup, usecase *usecase.TransactionUsecase) {
	handler := &TransactionHandler{usecase: usecase}

	api := r.Group("/transactions") 
	{
		api.POST("/", handler.Checkout)
	}
}

func (h *TransactionHandler) Checkout(c *gin.Context) {
	var trx models.Transaction

	if err := c.ShouldBindBodyWithJSON(&trx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format data salah" + err.Error()})
		return
	}

	if err := h.usecase.Checkout(&trx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal transaksi" + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "transaksi berhasil diproses!",
		"data": trx,
	})
}


