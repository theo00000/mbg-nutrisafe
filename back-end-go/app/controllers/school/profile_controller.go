package school

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
)

type UpdateSchoolProfileInput struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Lokasi string `json:"lokasi"`
	NPSN   string `json:"npsn"`
}

func GetSchoolProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User tidak ditemukan."})
	}

	var profile models.SchoolProfile
	config.DB.Where("user_id = ?", userID).First(&profile)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"phone":  user.Phone,
			"role":   user.RoleName,
			"lokasi": profile.Address,
			"npsn":   profile.NPSN,
			"grade":  profile.Grade,
		},
	})
}

func UpdateSchoolProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var input UpdateSchoolProfileInput
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

	var profile models.SchoolProfile
	config.DB.Where("user_id = ?", userID).FirstOrCreate(&profile, models.SchoolProfile{UserID: userID})
	if input.Lokasi != "" {
		profile.Address = input.Lokasi
	}
	if input.NPSN != "" {
		profile.NPSN = input.NPSN
	}
	config.DB.Save(&profile)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Profil sekolah berhasil diperbarui.",
		"data": fiber.Map{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"phone":  user.Phone,
			"lokasi": profile.Address,
			"npsn":   profile.NPSN,
		},
	})
}
