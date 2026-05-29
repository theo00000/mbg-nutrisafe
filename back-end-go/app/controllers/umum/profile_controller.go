package umum

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
)

func GetUmumProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Pengguna tidak ditemukan"})
	}

	var profile models.UmumProfile
	config.DB.Where("user_id = ?", userID).First(&profile)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"phone":  user.Phone,
			"role":   user.RoleName,
			"lokasi": profile.Lokasi,
		},
	})
}

func UpdateUmumProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var input struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Phone  string `json:"phone"`
		Lokasi string `json:"lokasi"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input tidak valid"})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Pengguna tidak ditemukan"})
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Phone != "" {
		user.Phone = input.Phone
	}
	config.DB.Save(&user)

	var profile models.UmumProfile
	config.DB.Where("user_id = ?", userID).FirstOrCreate(&profile, models.UmumProfile{UserID: userID})
	if input.Lokasi != "" {
		profile.Lokasi = input.Lokasi
	}
	config.DB.Save(&profile)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Profil berhasil diperbarui",
		"data": fiber.Map{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"phone":  user.Phone,
			"role":   user.RoleName,
			"lokasi": profile.Lokasi,
		},
	})
}
