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
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: file .env tidak ditemukan, menggunakan environment variabel dari server/sistem")
	}

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung dengan database: ", err)
	}

	DB = database
	log.Println("Berhasil terhubung ke database PostgreSQL!")

	err = DB.AutoMigrate(
		&models.Store{},
		&models.User{},
		&models.Product{},
		&models.Transaction{},
		&models.TransactionDetail{},
	)

	if err != nil {
		log.Fatal("Gagal melakukan migrasi database: ", err)
	}
	log.Println("Seluruh tabel berhasil dibuat / dimigrasi ke database!")
}