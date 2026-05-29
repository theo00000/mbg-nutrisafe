package admin

import (
	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
)

var validReportStatuses = map[string]bool{
	"belum_ditinjau":  true,
	"sedang_ditinjau": true,
	"selesai":         true,
}

func GetReports(c *fiber.Ctx) error {
	status := c.Query("status")
	if status != "" && !validReportStatuses[status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Filter status tidak valid. Pilihan: belum_ditinjau, sedang_ditinjau, selesai.",
		})
	}

	var reports []models.FoodReport
	query := config.DB.Order("created_at desc")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&reports).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data laporan.",
		})
	}

	if reports == nil {
		reports = []models.FoodReport{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   reports,
	})
}

func GetReportDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	var report models.FoodReport
	if err := config.DB.First(&report, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Laporan tidak ditemukan.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   report,
	})
}

func UpdateReportStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var input models.StatusUpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format request tidak valid.",
		})
	}

	if input.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Field status tidak boleh kosong.",
		})
	}

	if !validReportStatuses[input.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Status tidak valid. Pilihan: belum_ditinjau, sedang_ditinjau, selesai.",
		})
	}

	var report models.FoodReport
	if err := config.DB.First(&report, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Laporan tidak ditemukan.",
		})
	}

	if err := config.DB.Model(&report).Update("status", input.Status).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal memperbarui status laporan.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Status laporan berhasil diperbarui.",
		"data":    report,
	})
}
