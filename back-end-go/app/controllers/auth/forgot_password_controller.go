package auth

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"back-end/app/models"
	"back-end/config"
	"back-end/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

const resetCodeTTL = 15 * time.Minute
const resetMaxAttempts = 5

func generateResetCode() (string, error) {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	n := (uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])) % 1000000
	return fmt.Sprintf("%06d", n), nil
}

type ForgotPasswordRequestInput struct {
	Email string `json:"email"`
}

// resolveUserByEmail mencari user berdasarkan email yang diinput. Karena login email
// untuk role school/sppg di-generate sistem (xxx@sch.id / xxx@nutrisafe.id) sementara
// user mengingat email pribadinya, lookup dilakukan berurutan ke:
//   1. users.email                          (umum/admin: login pakai email asli)
//   2. school_profiles.contact_email        (school: email pribadi guru)
//   3. sppg_profiles.contact_email          (sppg: email pribadi PIC)
// Mengembalikan user record, deliveryEmail (alamat inbox real untuk kirim kode), dan flag found.
func resolveUserByEmail(input string) (models.User, string, bool) {
	email := strings.TrimSpace(strings.ToLower(input))
	if email == "" {
		return models.User{}, "", false
	}

	var user models.User
	if err := config.DB.Where("LOWER(email) = ?", email).First(&user).Error; err == nil {
		return user, user.Email, true
	}

	var schoolProfile models.SchoolProfile
	if err := config.DB.Where("LOWER(contact_email) = ?", email).First(&schoolProfile).Error; err == nil {
		if err := config.DB.First(&user, schoolProfile.UserID).Error; err == nil {
			return user, schoolProfile.ContactEmail, true
		}
	}

	var sppgProfile models.SppgProfile
	if err := config.DB.Where("LOWER(contact_email) = ?", email).First(&sppgProfile).Error; err == nil {
		if err := config.DB.First(&user, sppgProfile.UserID).Error; err == nil {
			return user, sppgProfile.ContactEmail, true
		}
	}

	return models.User{}, "", false
}

// RequestPasswordReset menerima email, generate kode 6-digit, simpan hash + kirim email.
// Selalu balas sukses untuk mencegah enumeration.
func RequestPasswordReset(c *fiber.Ctx) error {
	var input ForgotPasswordRequestInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format request tidak valid."})
	}

	inputEmail := strings.TrimSpace(strings.ToLower(input.Email))
	if inputEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Email wajib diisi."})
	}

	successResp := fiber.Map{
		"status":  "success",
		"message": "Jika email terdaftar, kode verifikasi akan dikirim ke email tersebut.",
	}

	user, deliveryEmail, found := resolveUserByEmail(inputEmail)
	if !found {
		return c.Status(fiber.StatusOK).JSON(successResp)
	}

	code, err := generateResetCode()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal membuat kode verifikasi."})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengenkripsi kode verifikasi."})
	}

	// PasswordReset.Email disimpan dengan canonical login email (users.email) supaya
	// verify step konsisten apapun email yang diinput user.
	canonicalEmail := strings.ToLower(user.Email)

	config.DB.Model(&models.PasswordReset{}).
		Where("email = ? AND used = ?", canonicalEmail, false).
		Update("used", true)

	reset := models.PasswordReset{
		Email:     canonicalEmail,
		CodeHash:  string(hash),
		ExpiresAt: time.Now().Add(resetCodeTTL),
	}
	if err := config.DB.Create(&reset).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan kode verifikasi."})
	}

	go func() {
		if err := utils.SendPasswordResetCode(deliveryEmail, code); err != nil {
			fmt.Printf("[mailer] gagal kirim reset code ke %s: %v\n", deliveryEmail, err)
		}
	}()

	return c.Status(fiber.StatusOK).JSON(successResp)
}

type ResetPasswordInput struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

func VerifyPasswordReset(c *fiber.Ctx) error {
	var input ResetPasswordInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format request tidak valid."})
	}

	inputEmail := strings.TrimSpace(strings.ToLower(input.Email))
	code := strings.TrimSpace(input.Code)
	if inputEmail == "" || code == "" || input.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Email, kode verifikasi, dan password baru wajib diisi."})
	}
	if len(input.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Password baru minimal 8 karakter."})
	}

	user, _, found := resolveUserByEmail(inputEmail)
	if !found {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Kode verifikasi tidak valid atau sudah kadaluarsa."})
	}
	canonicalEmail := strings.ToLower(user.Email)

	var reset models.PasswordReset
	if err := config.DB.
		Where("email = ? AND used = ?", canonicalEmail, false).
		Order("created_at desc").
		First(&reset).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Kode verifikasi tidak valid atau sudah kadaluarsa."})
	}

	if time.Now().After(reset.ExpiresAt) {
		config.DB.Model(&reset).Update("used", true)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Kode verifikasi sudah kadaluarsa. Mohon minta kode baru."})
	}

	if reset.Attempts >= resetMaxAttempts {
		config.DB.Model(&reset).Update("used", true)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Terlalu banyak percobaan. Mohon minta kode baru."})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(reset.CodeHash), []byte(code)); err != nil {
		config.DB.Model(&reset).Update("attempts", reset.Attempts+1)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Kode verifikasi tidak sesuai."})
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengenkripsi password baru."})
	}

	user.PasswordHash = string(newHash)
	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan password baru."})
	}

	config.DB.Model(&reset).Update("used", true)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Password berhasil diperbarui. Silakan login dengan password baru."})
}
