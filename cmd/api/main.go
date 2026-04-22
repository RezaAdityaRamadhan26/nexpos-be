package api

import (
	"net/http"
	"time"

	"nexpos-be/internal/config"
	deliveryHTTP "nexpos-be/internal/delivery/http"
	"nexpos-be/internal/delivery/http/middleware"
	"nexpos-be/internal/repository"
	"nexpos-be/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var app *gin.Engine

func init() {
	config.ConnectDB()

	app = gin.Default()

	app.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	userRepo := repository.NewUserRepository(config.DB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	deliveryHTTP.NewUserHandler(app, userUsecase)

	productRepo := repository.NewProductRepository(config.DB)
	productUsecase := usecase.NewProductUsecase(productRepo)

	transactionRepo := repository.NewTransactionRepository(config.DB)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo)

	protected := app.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	deliveryHTTP.NewProductHandler(protected, productUsecase)
	deliveryHTTP.NewTransactionHandler(protected, transactionUsecase)

	app.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "berhasil dari vercel!"})
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}