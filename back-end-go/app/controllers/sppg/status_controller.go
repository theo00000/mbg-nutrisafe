package sppg

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
)

func GetPartnershipStatus(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User tidak ditemukan."})
	}

	var profiles []models.SppgProfile
	config.DB.Where("user_id = ?", userID).Find(&profiles)

	var schoolCount int64
	config.DB.Model(&models.SchoolProfile{}).Where("sppg_user_id = ?", userID).Count(&schoolCount)

	type sppgItem struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		Status   string `json:"status"`
		Address  string `json:"address"`
		Capacity int    `json:"capacity"`
	}
	sppgList := make([]sppgItem, 0)
	for _, p := range profiles {
		sppgList = append(sppgList, sppgItem{
			ID:       p.ID,
			Name:     user.Name,
			Status:   p.Status,
			Address:  p.Address,
			Capacity: p.Capacity,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"is_partner":   len(profiles) > 0,
			"sppg_count":   len(profiles),
			"school_count": schoolCount,
			"sppg_list":    sppgList,
		},
	})
}
