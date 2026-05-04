package controllers

import (
	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
)

func GetPublicStats(c *fiber.Ctx) error {
	var schoolCount int64
	var spggCount int64
	var studentCount int64

	if err := config.DB.Model(&models.User{}).Where("role_name = ?", "school").Count(&schoolCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":       "error",
			"message":      "Gagal mengambil data jumlah sekolah",
			"error_detail": err.Error(),
		})
	}

	if err := config.DB.Model(&models.User{}).Where("role_name = ?", "spgg").Count(&spggCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":       "error",
			"message":      "Gagal mengambil data jumlah mitra SPPG",
			"error_detail": err.Error(),
		})
	}

	if err := config.DB.Model(&models.Student{}).Count(&studentCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":       "error",
			"message":      "Gagal mengambil data jumlah siswa",
			"error_detail": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"total_school":  schoolCount,
			"total_spgg":    spggCount,
			"total_student": studentCount,
		},
	})
}