package main

import (
	"nexpos-be/internal/config"
	deliveryHTTP "nexpos-be/internal/delivery/http" // alias biar ga tabrakan sama package bawaan go
	"nexpos-be/internal/delivery/http/middleware"
	"nexpos-be/internal/repository"
	"nexpos-be/internal/usecase"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// user
	userRepo := repository.NewUserRepository(config.DB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	deliveryHTTP.NewUserHandler(r, userUsecase)

	// product
	productRepo := repository.NewProductRepository(config.DB)
	productUsecase := usecase.NewProductUsecase(productRepo)

	// transaction
	transactionRepo := repository.NewTransactionRepository(config.DB)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo)
	
	// middleware
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	// protected
	deliveryHTTP.NewProductHandler(protected, productUsecase)
	deliveryHTTP.NewTransactionHandler(protected, transactionUsecase)

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong!"})
	})

	r.Run(":8080")
}