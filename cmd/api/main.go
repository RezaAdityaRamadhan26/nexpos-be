package main

import (
	"nexpos-be/internal/config"
	deliveryHTTP "nexpos-be/internal/delivery/http" // alias biar ga tabrakan sama package bawaan go
	"nexpos-be/internal/repository"
	"nexpos-be/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	r := gin.Default()

	userRepo := repository.NewUserRepository(config.DB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	deliveryHTTP.NewUserHandler(r, userUsecase)

	productRepo := repository.NewProductRepository(config.DB)
	productUsecase := usecase.NewProductUsecase(*productRepo)
	deliveryHTTP.NewProductHandler(r, productUsecase)

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong!"})
	})

	r.Run(":8080")
}