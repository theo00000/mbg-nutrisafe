package sppg

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
)

func GetAvailableSchools(c *fiber.Ctx) error {
	type SchoolItem struct {
		ID           uint   `json:"id"`
		Name         string `json:"name"`
		Grade        string `json:"grade"`
		StudentCount int    `json:"student_count"`
	}

	var profiles []models.SchoolProfile
	if err := config.DB.Where("sppg_user_id IS NULL").Find(&profiles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil data sekolah."})
	}

	result := make([]SchoolItem, 0)
	for _, sp := range profiles {
		var school models.User
		if config.DB.First(&school, sp.UserID).Error != nil {
			continue
		}
		result = append(result, SchoolItem{
			ID:           school.ID,
			Name:         school.Name,
			Grade:        sp.Grade,
			StudentCount: sp.StudentCount,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func AssignSchools(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var body struct {
		SchoolIDs []uint `json:"school_ids"`
	}
	if err := c.BodyParser(&body); err != nil || len(body.SchoolIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "school_ids tidak valid"})
	}

	if err := config.DB.Model(&models.SchoolProfile{}).
		Where("user_id IN ? AND sppg_user_id IS NULL", body.SchoolIDs).
		Update("sppg_user_id", userID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menetapkan sekolah."})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Sekolah berhasil ditambahkan."})
}

func GetsppgSchools(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var schoolProfiles []models.SchoolProfile
	if err := config.DB.Where("sppg_user_id = ?", userID).Find(&schoolProfiles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil data sekolah."})
	}

	type SchoolItem struct {
		ID           uint     `json:"id"`
		Name         string   `json:"name"`
		Address      string   `json:"address"`
		Grade        string   `json:"Grade"`
		TeacherCount int      `json:"teacher_count"`
		StudentCount int      `json:"student_count"`
		AllergyTypes []string `json:"allergy_types"`
	}

	totalStudents := 0
	result := make([]SchoolItem, 0)
	for _, sp := range schoolProfiles {
		var school models.User
		if config.DB.First(&school, sp.UserID).Error != nil {
			continue
		}

		var allergyTypes []string
		config.DB.Model(&models.StudentAllergy{}).
			Where("school_id = ?", sp.UserID).
			Distinct("allergy_type").
			Pluck("allergy_type", &allergyTypes)

		totalStudents += sp.StudentCount
		result = append(result, SchoolItem{
			ID:           school.ID,
			Name:         school.Name,
			Address:      sp.Address,
			Grade:        sp.Grade,
			TeacherCount: sp.TeacherCount,
			StudentCount: sp.StudentCount,
			AllergyTypes: allergyTypes,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"total_school":   len(result),
			"total_students": totalStudents,
			"schools":        result,
		},
	})
}
