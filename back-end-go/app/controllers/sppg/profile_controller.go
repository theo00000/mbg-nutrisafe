package sppg

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
)

type UpdatesppgProfileInput struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Lokasi          string `json:"lokasi"`
	InstitutionName string `json:"institution_name"`
}

func GetsppgProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User tidak ditemukan."})
	}

	var profile models.SppgProfile
	config.DB.Where("user_id = ?", userID).First(&profile)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":               user.ID,
			"name":             user.Name,
			"email":            user.Email,
			"phone":            user.Phone,
			"role":             user.RoleName,
			"lokasi":           profile.Address,
			"institution_name": profile.InstitutionName,
			"capacity":         profile.Capacity,
			"sppg_status":      profile.Status,
		},
	})
}

func UpdatesppgProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var input UpdatesppgProfileInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format request tidak valid."})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User tidak ditemukan."})
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" && input.Email != user.Email {
		var check models.User
		if config.DB.Where("email = ?", input.Email).First(&check).Error == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "error", "message": "Email sudah digunakan oleh pengguna lain."})
		}
		user.Email = input.Email
	}
	if input.Phone != "" && input.Phone != user.Phone {
		var check models.User
		if config.DB.Where("phone = ?", input.Phone).First(&check).Error == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "error", "message": "Nomor telepon sudah digunakan oleh pengguna lain."})
		}
		user.Phone = input.Phone
	}
	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan data pengguna."})
	}

	var profile models.SppgProfile
	config.DB.Where("user_id = ?", userID).FirstOrCreate(&profile, models.SppgProfile{UserID: userID})
	if input.Lokasi != "" {
		profile.Address = input.Lokasi
	}
	if input.InstitutionName != "" {
		profile.InstitutionName = input.InstitutionName
	}
	config.DB.Save(&profile)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Profil berhasil diperbarui.",
		"data": fiber.Map{
			"id":               user.ID,
			"name":             user.Name,
			"email":            user.Email,
			"phone":            user.Phone,
			"lokasi":           profile.Address,
			"institution_name": profile.InstitutionName,
		},
	})
}
