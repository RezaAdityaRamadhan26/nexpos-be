package main

import (
	"nexpos-be/internal/config"
	deliveryHTTP "nexpos-be/internal/delivery/http" // alias biar nama nya ga tabrakan sama package http bawaan go
	"nexpos-be/internal/repository"
	"nexpos-be/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	r := gin.Default()

	productRepo := repository.NewProductRepository(config.DB)

	ProductUsecase := usecase.NewProductUsecase(*productRepo)

	deliveryHTTP.NewProductHandler(r, ProductUsecase)

	r.GET("/api/ping", func (c *gin.Context){
		c.JSON(200, gin.H{"message": "pong!"})
	})

	r.Run(":8080")
}