package school

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type TeacherInput struct {
	Name     string `json:"name"`
	NIP      string `json:"nip"`
	Position string `json:"position"`
	Gender   string `json:"gender"`
	Address  string `json:"address"`
}

var validGenders = map[string]bool{
	"Laki-Laki": true,
	"Perempuan": true,
}

func GetTeachers(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var teachers []models.Teacher
	if err := config.DB.Where("school_id = ?", userID).Order("position asc, name asc").Find(&teachers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil data guru."})
	}
	if teachers == nil {
		teachers = []models.Teacher{}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": teachers})
}

func validateTeacher(input TeacherInput) (TeacherInput, *fiber.Map) {
	input.Name = strings.TrimSpace(input.Name)
	input.NIP = strings.TrimSpace(input.NIP)
	input.Position = strings.TrimSpace(input.Position)
	input.Gender = strings.TrimSpace(input.Gender)
	input.Address = strings.TrimSpace(input.Address)

	if input.Name == "" || input.NIP == "" || input.Position == "" || input.Gender == "" {
		return input, &fiber.Map{"status": "error", "message": "Nama, NIP, jabatan, dan jenis kelamin wajib diisi."}
	}
	if !validGenders[input.Gender] {
		return input, &fiber.Map{"status": "error", "message": "Jenis kelamin tidak valid. Pilihan: Laki-Laki, Perempuan."}
	}
	return input, nil
}

func AddTeacher(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var input TeacherInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format request tidak valid."})
	}
	input, errMap := validateTeacher(input)
	if errMap != nil {
		return c.Status(fiber.StatusBadRequest).JSON(*errMap)
	}

	teacher := models.Teacher{
		SchoolID: userID,
		Name:     input.Name,
		NIP:      input.NIP,
		Position: input.Position,
		Gender:   input.Gender,
		Address:  input.Address,
	}
	if err := config.DB.Create(&teacher).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan data guru."})
	}

	syncTeacherCount(userID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Data guru berhasil disimpan.",
		"data":    teacher,
	})
}

func syncTeacherCount(schoolID uint) {
	var count int64
	config.DB.Model(&models.Teacher{}).Where("school_id = ?", schoolID).Count(&count)
	config.DB.Model(&models.SchoolProfile{}).Where("user_id = ?", schoolID).Update("teacher_count", count)
}

func UpdateTeacher(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	id := c.Params("id")
	var teacher models.Teacher
	if err := config.DB.Where("id = ? AND school_id = ?", id, userID).First(&teacher).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Data guru tidak ditemukan."})
	}

	var input TeacherInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format request tidak valid."})
	}
	input, errMap := validateTeacher(input)
	if errMap != nil {
		return c.Status(fiber.StatusBadRequest).JSON(*errMap)
	}

	teacher.Name = input.Name
	teacher.NIP = input.NIP
	teacher.Position = input.Position
	teacher.Gender = input.Gender
	teacher.Address = input.Address
	if err := config.DB.Save(&teacher).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal memperbarui data guru."})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Data guru berhasil diperbarui.",
		"data":    teacher,
	})
}

func DeleteTeacher(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	id := c.Params("id")
	res := config.DB.Where("id = ? AND school_id = ?", id, userID).Delete(&models.Teacher{})
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menghapus data guru."})
	}
	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Data guru tidak ditemukan."})
	}

	syncTeacherCount(userID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Data guru berhasil dihapus.",
	})
}
