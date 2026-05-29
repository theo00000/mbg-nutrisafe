package sppg

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
)

func GetAllergyDashboard(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var schoolProfiles []models.SchoolProfile
	config.DB.Where("sppg_user_id = ?", userID).Find(&schoolProfiles)

	schoolIDs := make([]uint, 0)
	for _, sp := range schoolProfiles {
		schoolIDs = append(schoolIDs, sp.UserID)
	}

	if len(schoolIDs) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"total_allergy_types":    0,
				"total_allergy_students": 0,
				"allergy_breakdown":      []interface{}{},
				"schools":                []interface{}{},
			},
		})
	}

	type AllergyCount struct {
		AllergyType string `json:"allergy_type"`
		Count       int    `json:"count"`
	}
	var allergyBreakdown []AllergyCount
	config.DB.Model(&models.StudentAllergy{}).
		Select("allergy_type, count(*) as count").
		Where("school_id IN ?", schoolIDs).
		Group("allergy_type").
		Scan(&allergyBreakdown)

	var totalAllergyStudents int64
	config.DB.Model(&models.StudentAllergy{}).
		Where("school_id IN ?", schoolIDs).
		Count(&totalAllergyStudents)

	type SchoolAllergyItem struct {
		ID            uint   `json:"id"`
		Name          string `json:"name"`
		TotalStudents int    `json:"total_students"`
		AllergyCount  int64  `json:"allergy_count"`
	}
	schoolItems := make([]SchoolAllergyItem, 0)
	for _, sp := range schoolProfiles {
		var school models.User
		if config.DB.First(&school, sp.UserID).Error != nil {
			continue
		}
		var count int64
		config.DB.Model(&models.StudentAllergy{}).Where("school_id = ?", sp.UserID).Count(&count)
		schoolItems = append(schoolItems, SchoolAllergyItem{
			ID:            school.ID,
			Name:          school.Name,
			TotalStudents: sp.StudentCount,
			AllergyCount:  count,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"total_allergy_types":    len(allergyBreakdown),
			"total_allergy_students": totalAllergyStudents,
			"allergy_breakdown":      allergyBreakdown,
			"schools":                schoolItems,
		},
	})
}
