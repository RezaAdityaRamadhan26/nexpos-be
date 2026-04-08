package http

import (
	"net/http"
	"nexpos-be/internal/models"
	"nexpos-be/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	usecase *usecase.ProductUsecase
}

func NewProductHandler(r *gin.RouterGroup, usecase *usecase.ProductUsecase) {
	handler := &ProductHandler{usecase: usecase}

	api := r.Group("/products")	
	{
		api.POST("/", handler.Create)
		api.GET("/", handler.GetAll)
		api.PUT("/:id", handler.Update)
		api.DELETE("/:id", handler.Delete)
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

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {"error": "ID Produk tidak valid"})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.UpdateProducts(id, &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengupdate produk: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "produk berhasil diupdate", "data": product})
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID produk tidak valid!"})
		return
	}

	if err := h.usecase.DeleteProducts(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "produk berhasil dihapus"})
}

