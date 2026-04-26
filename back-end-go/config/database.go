package config

import (
	"fmt"
	"log"
	"os"

	"back-end/app/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, pastikan env tersedia di system")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database NutriSafe!\n", err)
	}

	fmt.Println("Koneksi ke PostgreSQL berhasil!")

	fmt.Println("Sedang menyinkronkan tabel database...")
	err = database.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Menu{},
		&models.Ingredient{},
		&models.Student{},
		&models.Allergy{},
		&models.FoodReport{},
		&models.DailyMenuPlan{},
		&models.AllergyAlternativeSuggestion{},
	)

	if err != nil {
		log.Fatal("Proses migrasi tabel gagal:", err)
	}

	fmt.Println("Migrasi tabel berhasil diselesaikan!")

	SeedRoles(database)

	DB = database
}

func SeedRoles(db *gorm.DB) {
	roles := []models.Role{
		{Name: "school", Description: "Akun representasi pihak sekolah"},
		{Name: "spgg", Description: "Akun untuk Ahli Gizi (SPGG)"},
		{Name: "umum", Description: "Akun untuk orang tua atau pengguna umum"},
	}

	for _, role := range roles {
		db.Where("name = ?", role.Name).FirstOrCreate(&role)
	}
	
	fmt.Println("Pengecekan dan pengisian 3 role dasar berhasil!")
}