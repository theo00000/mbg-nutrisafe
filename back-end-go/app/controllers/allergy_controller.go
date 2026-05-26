package controllers

import (
	"fmt"
	"os"
	"strings"

	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func getUserIDFromToken(c *fiber.Ctx) (uint, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("token tidak ditemukan")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, fmt.Errorf("format token tidak valid")
	}
	token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("token tidak valid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("gagal membaca claims")
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("id tidak valid")
	}
	return uint(userIDFloat), nil
}

type AllergyInput struct {
	StudentName    string `json:"student_name"`
	ClassName      string `json:"class_name"`
	AllergyType    string `json:"allergy_type"`
	Description    string `json:"description"`
	Severity       string `json:"severity"`
	ActionRequired string `json:"action_required"`
}

type SchoolInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name_school"`
}

type AllergyResponse struct {
	ID             uint       `json:"id"`
	School         SchoolInfo `json:"school"` 
	StudentName    string     `json:"student_name"`
	ClassName      string     `json:"class_name"`
	AllergyType    string     `json:"allergy_type"`
	Description    string     `json:"description"`
	Severity       string     `json:"severity"`
	ActionRequired string     `json:"action_required"`
}

func CreateAllergy(c *fiber.Ctx) error {
	schoolID, err := getUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var input AllergyInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format request tidak valid."})
	}

	if input.StudentName == "" || input.ClassName == "" || input.AllergyType == "" || input.Severity == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Nama Siswa, Kelas, Jenis Alergi, dan Tingkat Keparahan wajib diisi!",
		})
	}

	allergy := models.StudentAllergy{
		SchoolID:       schoolID, 
		StudentName:    input.StudentName,
		ClassName:      input.ClassName,
		AllergyType:    input.AllergyType,
		Description:    input.Description,
		Severity:       input.Severity,
		ActionRequired: input.ActionRequired,
	}

	if err := config.DB.Create(&allergy).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan data."})
	}

	var school models.User
	config.DB.Select("id, name").First(&school, schoolID)

	response := AllergyResponse{
		ID:             allergy.ID,
		School:         SchoolInfo{ID: school.ID, Name: school.Name},
		StudentName:    allergy.StudentName,
		ClassName:      allergy.ClassName,
		AllergyType:    allergy.AllergyType,
		Description:    allergy.Description,
		Severity:       allergy.Severity,
		ActionRequired: allergy.ActionRequired,
	}

	return c.Status(fiber.StatusCreated).JSON(struct {
		Status  string          `json:"status"`
		Message string          `json:"message"`
		Data    AllergyResponse `json:"data"`
	}{
		Status:  "success",
		Message: "Data alergi siswa berhasil disimpan.",
		Data:    response,
	})
}

func GetAllergies(c *fiber.Ctx) error {
	schoolID, err := getUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var school models.User
	if err := config.DB.Select("id, name").First(&school, schoolID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal memuat data sekolah."})
	}

	var allergies []models.StudentAllergy
	if err := config.DB.Where("school_id = ?", schoolID).Order("created_at desc").Find(&allergies).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil data."})
	}

	var responseData []AllergyResponse
	
	if len(allergies) == 0 {
		responseData = []AllergyResponse{}
	}

	for _, a := range allergies {
		responseData = append(responseData, AllergyResponse{
			ID: a.ID,
			School: SchoolInfo{
				ID:   school.ID,
				Name: school.Name,
			},
			StudentName:    a.StudentName,
			ClassName:      a.ClassName,
			AllergyType:    a.AllergyType,
			Description:    a.Description,
			Severity:       a.Severity,
			ActionRequired: a.ActionRequired,
		})
	}

	return c.Status(fiber.StatusOK).JSON(struct {
		Status string            `json:"status"`
		Data   []AllergyResponse `json:"data"`
	}{
		Status: "success",
		Data:   responseData,
	})
}