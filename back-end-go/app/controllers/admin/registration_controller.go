package admin

import (
	"fmt"
	"strings"

	"back-end/app/models"
	"back-end/config"
	"back-end/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var validRegistrationStatuses = map[string]bool{
	"pending":        true,
	"belum_ditinjau": true,
	"verifikasi":     true,
	"disetujui":      true,
}

func GetRegistrations(c *fiber.Ctx) error {
	status := c.Query("status")
	if status != "" && !validRegistrationStatuses[status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Filter status tidak valid. Pilihan: pending, belum_ditinjau, verifikasi, disetujui.",
		})
	}

	var registrations []models.SppgRegistration
	query := config.DB.Order("created_at desc")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&registrations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data pendaftaran.",
		})
	}

	if registrations == nil {
		registrations = []models.SppgRegistration{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   registrations,
	})
}

func GetRegistrationDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	var registration models.SppgRegistration
	if err := config.DB.First(&registration, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Data pendaftaran tidak ditemukan.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   registration,
	})
}

func UpdateRegistrationStatus(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ID pendaftaran wajib diisi pada URL.",
		})
	}

	var input models.StatusUpdateInput
	if err := c.BodyParser(&input); err != nil || strings.TrimSpace(input.Status) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Request tidak valid, pilih: ['belum_ditinjau', 'verifikasi', 'disetujui'].",
		})
	}

	input.Status = strings.TrimSpace(input.Status)
	if !validRegistrationStatuses[input.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Status '" + input.Status + "' tidak valid. Pilihan yang tersedia: ['belum_ditinjau', 'verifikasi', 'disetujui'].",
		})
	}

	var registration models.SppgRegistration
	if err := config.DB.First(&registration, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Data pendaftaran dengan ID " + id + " tidak ditemukan.",
		})
	}

	if err := config.DB.Model(&registration).Update("status", input.Status).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal memperbarui status pendaftaran.",
		})
	}

	if input.Status == "disetujui" {
		_, user, err := createsppgUserFromRegistration(registration)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Status disetujui, tapi gagal membuat akun sppg: " + err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Pendaftaran disetujui. Akun sppg berhasil dibuat. Kredensial login telah dikirim ke email PIC.",
			"data": fiber.Map{
				"registration": registration,
				"user_id":      user.ID,
				"email":        user.Email,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Status pendaftaran berhasil diperbarui.",
		"data":    registration,
	})
}

func createsppgUserFromRegistration(reg models.SppgRegistration) (string, models.User, error) {
	// generate login email unik dari nama SPPG
	slug := utils.GenerateSlug(reg.SppgName)
	loginEmail := slug + "@sppg.id"
	counter := 2
	for {
		var existing models.User
		if config.DB.Where("email = ?", loginEmail).First(&existing).Error != nil {
			break
		}
		loginEmail = fmt.Sprintf("%s.%d@sppg.id", slug, counter)
		counter++
	}

	tempPassword := utils.GeneratePasswordFromName(reg.SppgName)
	hash, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", models.User{}, err
	}

	user := models.User{
		Name:         reg.SppgName,
		Email:        loginEmail,
		Phone:        reg.Phone,
		PasswordHash: string(hash),
		RoleName:     "sppg",
	}
	if err := config.DB.Create(&user).Error; err != nil {
		return "", models.User{}, err
	}

	profile := models.SppgProfile{
		UserID:          user.ID,
		InstitutionName: reg.InstitutionName,
		Capacity:        reg.ProductionCapacity,
		Address:         reg.SppgAddress,
		ContactEmail:    reg.Email,
		Status:          "active",
	}
	if err := config.DB.Create(&profile).Error; err != nil {
		config.DB.Delete(&user)
		return "", models.User{}, err
	}

	go func() {
		if err := utils.SendSppgCredentials(reg.Email, reg.SppgName, loginEmail, tempPassword); err != nil {
			fmt.Printf("[mailer] gagal kirim email ke %s: %v\n", reg.Email, err)
		}
	}()

	return tempPassword, user, nil
}
