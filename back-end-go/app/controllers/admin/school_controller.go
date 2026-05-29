package admin

import (
	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
)

type schoolItem struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Grade        string `json:"grade"`
	Address      string `json:"address"`
	SppgName     string `json:"sppg_name"`
	TeacherCount int    `json:"teacher_count"`
	StudentCount int    `json:"student_count"`
	Status       string `json:"status"`
}

func buildSchoolItem(school models.User, profile models.SchoolProfile) schoolItem {
	sppgName := ""
	if profile.SppgUserID != nil {
		var sppg models.User
		if config.DB.First(&sppg, *profile.SppgUserID).Error == nil {
			sppgName = sppg.Name
		}
	}
	return schoolItem{
		ID:           school.ID,
		Name:         school.Name,
		Email:        school.Email,
		Grade:        profile.Grade,
		Address:      profile.Address,
		SppgName:     sppgName,
		TeacherCount: profile.TeacherCount,
		StudentCount: profile.StudentCount,
		Status:       profile.Status,
	}
}

func GetSchools(c *fiber.Ctx) error {
	var schools []models.User
	if err := config.DB.Where("role_name = ?", "school").Find(&schools).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data sekolah.",
		})
	}

	result := make([]schoolItem, 0)
	for _, school := range schools {
		var profile models.SchoolProfile
		config.DB.Where("user_id = ?", school.ID).First(&profile)
		result = append(result, buildSchoolItem(school, profile))
	}

	return c.Status(fiber.StatusOK).JSON(struct {
		Status string       `json:"status"`
		Data   []schoolItem `json:"data"`
	}{"success", result})
}

func GetSchoolDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	var school models.User
	if err := config.DB.Where("id = ? AND role_name = ?", id, "school").First(&school).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Sekolah tidak ditemukan.",
		})
	}

	var profile models.SchoolProfile
	config.DB.Where("user_id = ?", school.ID).First(&profile)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   buildSchoolItem(school, profile),
	})
}
