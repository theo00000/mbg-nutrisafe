package controllers

import (
	"os"
	"time"

	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleName string `json:"role_name"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	var input RegisterInput

	// Worst-case 1: Format JSON rusak atau tipe data salah
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format request tidak valid. Pastikan mengirim JSON yang benar.",
		})
	}

	// Worst-case 2: Ada field yang sengaja dikosongkan
	if input.Name == "" || input.Email == "" || input.Password == "" || input.RoleName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Semua kolom (name, email, password, role_name) wajib diisi!",
		})
	}

	// Worst-case 3: Email sudah terdaftar
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{ // 409 Conflict
			"status":  "error",
			"message": "Email ini sudah terdaftar. Silakan gunakan email lain atau langsung login.",
		})
	}

	// Worst-case 4: Role typo atau dikirim dengan role yang lain
	var role models.Role
	if err := config.DB.Where("name = ?", input.RoleName).First(&role).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Role tidak valid! Pilihan yang tersedia: school, spgg, umum.",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengenkripsi kata sandi.",
		})
	}

	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
	}

	// Worst-case 5: Database down
	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan data pengguna ke server. Silakan coba lagi nanti.",
		})
	}

	config.DB.Model(&user).Association("Roles").Append(&role)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Registrasi berhasil",
		"data": fiber.Map{
			"user_id": user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"role":    input.RoleName,
		},
	})
}

func Login(c *fiber.Ctx) error {
	var input LoginInput

	// Worst-case 1: JSON rusak
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format request tidak valid.",
		})
	}

	// Worst-case 2: Field kosong
	if input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Email dan kata sandi tidak boleh kosong!",
		})
	}

	var user models.User
	if err := config.DB.Preload("Roles").Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Email atau kata sandi salah.",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Email atau kata sandi salah.",
		})
	}

	userRole := "umum" 
	if len(user.Roles) > 0 {
		userRole = user.Roles[0].Name
	}

	// Worst-case 3: Kunci JWT tidak ada di .env
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Terjadi kesalahan internal server (Konfigurasi keamanan tidak lengkap).",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    userRole,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menerbitkan token sesi masuk.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login berhasil",
		"token":   tokenString,
		"role":    userRole,
	})
}