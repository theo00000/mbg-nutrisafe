package school

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
)

type ClassInput struct {
	Name string `json:"name"`
}

func GetClasses(c *fiber.Ctx) error {
	schoolID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var classes []models.SchoolClass
	config.DB.Where("school_id = ?", schoolID).Order("name asc").Find(&classes)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   classes,
	})
}

func AddClass(c *fiber.Ctx) error {
	schoolID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var input ClassInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format request tidak valid."})
	}

	if input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Nama kelas wajib diisi."})
	}

	var existing models.SchoolClass
	if config.DB.Where("school_id = ? AND name = ?", schoolID, input.Name).First(&existing).Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "error", "message": "Kelas dengan nama tersebut sudah ada."})
	}

	class := models.SchoolClass{
		SchoolID: schoolID,
		Name:     input.Name,
	}

	if err := config.DB.Create(&class).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan kelas."})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Kelas berhasil ditambahkan.",
		"data":    class,
	})
}

func DeleteClass(c *fiber.Ctx) error {
	schoolID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	id := c.Params("id")
	var class models.SchoolClass
	if err := config.DB.Where("id = ? AND school_id = ?", id, schoolID).First(&class).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Kelas tidak ditemukan."})
	}

	config.DB.Delete(&class)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Kelas berhasil dihapus."})
}
