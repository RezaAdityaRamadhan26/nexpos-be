package http

import (
	"net/http"
	"nexpos-be/internal/models"
	"nexpos-be/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	usecase *usecase.ProductUsecase
}

func NewProductHandler(r *gin.Engine, usecase *usecase.ProductUsecase) {
	handler := &ProductHandler{usecase: usecase}

	api := r.Group("/api/products")
	{
		api.POST("/", handler.Create)
		api.GET("/", handler.GetAll)
	}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product) ; err != nil {
		c.JSON(http.StatusBadRequest , gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.CreateProduct(&product) ; err != nil {
		c.JSON(http.StatusInternalServerError , gin.H{"error": "gagal menyimpan produk:" + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "produk berhasil ditambahkan",
		"data": product,
	})
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.usecase.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data produk" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}
