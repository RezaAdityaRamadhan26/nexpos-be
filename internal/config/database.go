package config

import (
	"fmt"
	"log"
	"os"

	"nexpos-be/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err:= godotenv.Load()
	if err!= nil {
		log.Println("peringatan! file .env tidak ditemukan, menggunakan environment default")
	}

	dsn := 	fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("gagal terhubung dengan database", err)
	}

	DB = database
	log.Println("terhubung ke database")

	err = DB.AutoMigrate(
		&models.Store{},
		&models.User{},
		&models.Product{},
		&models.Transaction{},
		&models.TransactionDetail{},
	)

	if err != nil {
		log.Fatal("gagal melakukan migrasi database", err)
	}
	log.Println("seluruh table berhasul dibuat / di migrasi ke database!")
}