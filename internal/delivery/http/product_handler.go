package http

import (
	"net/http"
	"nexpos-be/internal/models"
	"nexpos-be/internal/usecase"
	"strconv"
	"nexpos-be/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	usecase *usecase.ProductUsecase
}

func NewProductHandler(r *gin.RouterGroup, usecase *usecase.ProductUsecase) {
	handler := &ProductHandler{usecase: usecase}

	api := r.Group("/products")
	{
		api.POST("/", middleware.AuthMiddleware(), middleware.OwnerOnly(), handler.Create)
		api.GET("/", middleware.AuthMiddleware(), handler.GetAll)
		api.PUT("/:id", middleware.AuthMiddleware(), middleware.OwnerOnly(), handler.Update)
		api.DELETE("/:id", middleware.AuthMiddleware(), middleware.OwnerOnly(), handler.Delete)
	}
}
func getStoreIDFromContext(c *gin.Context) (uint, error) {
	storeID, exists := c.Get("store.id")
	if !exists {
		return 0, http.ErrNoCookie 
	}

	if idFloat, ok := storeID.(float64); ok {	
		return uint(idFloat), nil
	}
	return 0, http.ErrNotSupported
}

func (h *ProductHandler) Create(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tidak terautentikasi atau format token salah"})
		return
	}
	product.StoreID = storeID

	if err := h.usecase.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan produk: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "produk berhasil ditambahkan",
		"data":    product,
	})
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tidak terautentikasi atau format token salah"})
		return
	}

	products, errUsecase := h.usecase.GetAllProducts(storeID)
	if errUsecase != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data produk: " + errUsecase.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Produk tidak valid"})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storeID, errCtx := getStoreIDFromContext(c)
	if errCtx != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tidak terautentikasi atau format token salah"})
		return
	}

	if err := h.usecase.UpdateProducts(id, storeID, &product); err != nil {
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

	storeID, errCtx := getStoreIDFromContext(c)
	if errCtx != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tidak terautentikasi atau format token salah"})
		return
	}

	if err := h.usecase.DeleteProducts(id, storeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "produk berhasil dihapus"})
}