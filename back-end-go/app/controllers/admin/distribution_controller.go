package admin

import (
	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
)

var validDistributionStatuses = map[string]bool{
	"dalam_pengiriman": true,
	"diterima":         true,
}

func GetDistributions(c *fiber.Ctx) error {
	dateFilter := c.Query("date")
	distributions := make([]models.FoodDistribution, 0)
	query := config.DB.Order("created_at desc")
	if dateFilter != "" {
		query = query.Where("dist_date = ?", dateFilter)
	}
	if err := query.Find(&distributions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data distribusi.",
		})
	}

	var total, inDelivery, received int64
	config.DB.Model(&models.FoodDistribution{}).Count(&total)
	config.DB.Model(&models.FoodDistribution{}).Where("status = ?", "dalam_pengiriman").Count(&inDelivery)
	config.DB.Model(&models.FoodDistribution{}).Where("status = ?", "diterima").Count(&received)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"summary": fiber.Map{
				"total":       total,
				"in_delivery": inDelivery,
				"received":    received,
			},
			"school": distributions,
		},
	})
}

func UpdateDistributionStatus(c *fiber.Ctx) error {
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

	if !validDistributionStatuses[input.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Status tidak valid. Pilihan: dalam_pengiriman, diterima.",
		})
	}

	var distribution models.FoodDistribution
	if err := config.DB.First(&distribution, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Data distribusi tidak ditemukan.",
		})
	}

	if err := config.DB.Model(&distribution).Update("status", input.Status).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal memperbarui status distribusi.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Status distribusi berhasil diperbarui.",
		"data":    distribution,
	})
}
