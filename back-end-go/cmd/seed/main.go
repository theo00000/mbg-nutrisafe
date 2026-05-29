package main

import (
	"fmt"
	"log"
	"os"

	"back-end/app/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("File .env tidak ditemukan, pakai env system")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	email := "admin@nutrisafe.com"
	password := "admin123"
	name := "Super Admin"
	phone := "081234567890"

	var existing models.User
	if db.Where("email = ?", email).First(&existing).Error == nil {
		fmt.Printf("Admin dengan email '%s' sudah ada (ID: %d), skip.\n", email, existing.ID)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Gagal hash password:", err)
	}

	admin := models.User{
		Name:         name,
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hash),
		RoleName:     "admin",
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Fatal("Gagal membuat akun admin:", err)
	}

	fmt.Println("Akun admin berhasil dibuat!")
	fmt.Printf("  Email    : %s\n", email)
	fmt.Printf("  Password : %s\n", password)
	fmt.Printf("  Role     : admin\n")
}
