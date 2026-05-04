package controllers

import (
	"fmt"
	"os"
	"time"
	"strings"

	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	RoleName string `json:"role_name"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
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

	// Worst-case 3: Email atau No Telp sudah terdaftar
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{ // 409 Conflict
			"status":  "error",
			"message": "Email ini sudah terdaftar. Silakan gunakan email lain atau langsung login.",
		})
	}

	if err := config.DB.Where("phone = ?", input.Phone).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{ 
			"status":  "error",
			"message": "Nomor telepon ini sudah terdaftar. Silakan gunakan nomor lain.",
		})
	}

	// Worst-case 4: Role typo atau tidak valid (Validasi manual agar lebih cepat tanpa hit database)
	if input.RoleName != "school" && input.RoleName != "spgg" && input.RoleName != "umum" {
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
		Phone:        input.Phone,
		PasswordHash: string(hashedPassword),
		RoleName:     input.RoleName,
	}

	// Worst-case 5: Database down
	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan data pengguna ke server. Silakan coba lagi nanti.",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Registrasi berhasil",
		"data": fiber.Map{
			"user_id": user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"phone":   user.Phone,
			"role":    user.RoleName,
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
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
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

	userRole := user.RoleName

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

func Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal logout: Tidak ada token yang dikirim.",
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Format token tidak valid. Gunakan format 'Bearer <token>'.",
		})
	}

	tokenString := parts[1]

	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode enkripsi tidak valid")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Token tidak valid atau sudah kadaluarsa.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Logout berhasil.",
	})
}

func GetProfile(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Akses ditolak: Token tidak ditemukan.",
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Format token tidak valid.",
		})
	}
	tokenString := parts[1]

	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode enkripsi tidak valid")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Token tidak valid atau sudah kadaluarsa.",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal membaca data dari dalam token.",
		})
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Format ID user tidak valid di dalam token.",
		})
	}
	userID := uint(userIDFloat)

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User tidak ditemukan di database.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(struct {
		Status string      `json:"status"`
		Data   models.User `json:"data"`
	}{
		Status: "success",
		Data:   user,
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Akses ditolak: Token tidak ditemukan."})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Format token tidak valid."})
	}

	tokenString := parts[1]
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode enkripsi tidak valid")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Token tidak valid atau sudah kadaluarsa."})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal membaca data dari dalam token."})
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Format ID user tidak valid di dalam token."})
	}
	userID := uint(userIDFloat)


	var input UpdateProfileInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format request tidak valid.",
		})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User tidak ditemukan.",
		})
	}

	if input.Email != "" && input.Email != user.Email {
		var checkEmail models.User
		if err := config.DB.Where("email = ?", input.Email).First(&checkEmail).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status":  "error",
				"message": "Email ini sudah digunakan oleh pengguna lain.",
			})
		}
		user.Email = input.Email 
	}

	if input.Phone != "" && input.Phone != user.Phone {
		var checkPhone models.User
		if err := config.DB.Where("phone = ?", input.Phone).First(&checkPhone).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status":  "error",
				"message": "Nomor telepon ini sudah digunakan oleh pengguna lain.",
			})
		}
		user.Phone = input.Phone
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Phone != "" {
		user.Phone = input.Phone
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan pembaruan profil.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    models.User `json:"data"`
	}{
		Status:  "success",
		Message: "Profil berhasil diperbarui.",
		Data:    user,
	})
}