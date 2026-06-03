package admin

import (
	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
)

// Status di tabel DeliveryReport: "pending" / "diterima".
// FE admin lama menggunakan kosakata "dalam_pengiriman" / "diterima";
// kami terjemahkan otomatis di kedua arah supaya kontrak lama FE tidak berubah.

var distStatusFromAdmin = map[string]string{
	"dalam_pengiriman": "pending",
	"diterima":         "diterima",
}

func adminStatusOf(s string) string {
	if s == "pending" {
		return "dalam_pengiriman"
	}
	return s
}

type adminDistributionRow struct {
	ID       uint   `json:"id"`
	SchoolID uint   `json:"school_id"`
	SppgID   uint   `json:"sppg_id"`
	DistDate string `json:"dist_date"`
	Portions int    `json:"portions"`
	Status   string `json:"status"`
}

func GetDistributions(c *fiber.Ctx) error {
	dateFilter := c.Query("date")

	query := config.DB.Model(&models.DeliveryReport{}).Order("delivery_date desc, created_at desc")
	if dateFilter != "" {
		query = query.Where("delivery_date = ?", dateFilter)
	}

	var reports []models.DeliveryReport
	if err := query.Find(&reports).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data distribusi.",
		})
	}

	rows := make([]adminDistributionRow, 0, len(reports))
	for _, r := range reports {
		rows = append(rows, adminDistributionRow{
			ID:       r.ID,
			SchoolID: r.SchoolID,
			SppgID:   r.SppgID,
			DistDate: r.DeliveryDate.Format("2006-01-02"),
			Portions: r.TotalPortions,
			Status:   adminStatusOf(r.Status),
		})
	}

	var total, inDelivery, received int64
	config.DB.Model(&models.DeliveryReport{}).Count(&total)
	config.DB.Model(&models.DeliveryReport{}).Where("status = ?", "pending").Count(&inDelivery)
	config.DB.Model(&models.DeliveryReport{}).Where("status = ?", "diterima").Count(&received)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"summary": fiber.Map{
				"total":       total,
				"in_delivery": inDelivery,
				"received":    received,
			},
			"school": rows,
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

	mapped, ok := distStatusFromAdmin[input.Status]
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Status tidak valid. Pilihan: dalam_pengiriman, diterima.",
		})
	}

	var report models.DeliveryReport
	if err := config.DB.First(&report, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Data distribusi tidak ditemukan.",
		})
	}

	if err := config.DB.Model(&report).Update("status", mapped).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal memperbarui status distribusi.",
		})
	}
	report.Status = mapped

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Status distribusi berhasil diperbarui.",
		"data": adminDistributionRow{
			ID:       report.ID,
			SchoolID: report.SchoolID,
			SppgID:   report.SppgID,
			DistDate: report.DeliveryDate.Format("2006-01-02"),
			Portions: report.TotalPortions,
			Status:   adminStatusOf(report.Status),
		},
	})
}
