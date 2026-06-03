package school

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

var onlyDigits = regexp.MustCompile(`^\d+$`)

type StudentInput struct {
	Name    string `json:"name"`
	NISN    string `json:"nisn"`
	Class   string `json:"class"`
	Gender  string `json:"gender"`
	Address string `json:"address"`
}

func AddStudent(c *fiber.Ctx) error {
	schoolID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var input StudentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format request tidak valid."})
	}

	if input.Name == "" || input.NISN == "" || input.Class == "" || input.Gender == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Kolom name, nisn, class, dan gender wajib diisi.",
		})
	}

	if len(input.NISN) != 10 || !onlyDigits.MatchString(input.NISN) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "NISN harus terdiri dari 10 digit angka."})
	}

	if input.Gender != "Laki-laki" && input.Gender != "Perempuan" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Jenis kelamin harus 'Laki-laki' atau 'Perempuan'."})
	}

	var existing models.Student
	if config.DB.Where("nisn = ?", input.NISN).First(&existing).Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "error", "message": "NISN ini sudah terdaftar."})
	}

	student := models.Student{
		SchoolID: schoolID,
		Name:     input.Name,
		NISN:     input.NISN,
		Class:    input.Class,
		Gender:   input.Gender,
		Address:  input.Address,
	}

	if err := config.DB.Create(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan data siswa."})
	}

	syncStudentCount(schoolID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Data siswa berhasil disimpan.",
		"data":    student,
	})
}

func syncStudentCount(schoolID uint) {
	var count int64
	config.DB.Model(&models.Student{}).Where("school_id = ?", schoolID).Count(&count)
	config.DB.Model(&models.SchoolProfile{}).Where("user_id = ?", schoolID).Update("student_count", count)
}

func UpdateStudent(c *fiber.Ctx) error {
	schoolID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	id := c.Params("id")
	var student models.Student
	if err := config.DB.Where("id = ? AND school_id = ?", id, schoolID).First(&student).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Data siswa tidak ditemukan."})
	}

	if student.AccountGenerated {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Data siswa tidak dapat diubah karena akun sudah dibuat.",
		})
	}

	var input StudentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format request tidak valid."})
	}

	if input.Gender != "" && input.Gender != "Laki-laki" && input.Gender != "Perempuan" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Jenis kelamin harus 'Laki-laki' atau 'Perempuan'."})
	}

	if input.Name != "" {
		student.Name = input.Name
	}
	if input.Class != "" {
		student.Class = input.Class
	}
	if input.Gender != "" {
		student.Gender = input.Gender
	}
	if input.Address != "" {
		student.Address = input.Address
	}

	if err := config.DB.Save(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal memperbarui data siswa."})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Data siswa berhasil diperbarui.",
		"data":    student,
	})
}

func GetStudents(c *fiber.Ctx) error {
	schoolID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var students []models.Student
	config.DB.Where("school_id = ?", schoolID).Find(&students)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   students,
	})
}

func DeleteStudent(c *fiber.Ctx) error {
	schoolID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	id := c.Params("id")
	var student models.Student
	if err := config.DB.Where("id = ? AND school_id = ?", id, schoolID).First(&student).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Data siswa tidak ditemukan."})
	}

	if student.AccountGenerated {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Data siswa tidak dapat dihapus karena akun sudah dibuat.",
		})
	}

	config.DB.Delete(&student)
	syncStudentCount(schoolID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Data siswa berhasil dihapus."})
}
