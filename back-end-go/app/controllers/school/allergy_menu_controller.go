package school

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetAllergyMenu mengembalikan menu makanan khusus siswa alergi dari laporan pengiriman sppg.
// Dikelompokkan per jenis alergi, disertai jumlah siswa yang terdampak.
func GetAllergyMenu(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	// Default tanggal = hari ini
	dateStr := c.Query("date")
	var targetDate time.Time
	if dateStr != "" {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format tanggal tidak valid. Gunakan YYYY-MM-DD."})
		}
	} else {
		targetDate = time.Now().Truncate(24 * time.Hour)
	}

	// Cari delivery report untuk sekolah ini pada tanggal yang diminta
	var report models.DeliveryReport
	err = config.DB.
		Preload("Items").
		Where("school_id = ? AND delivery_date = ?", userID, targetDate.Format("2006-01-02")).
		First(&report).Error
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"date":    targetDate.Format("2006-01-02"),
			"data":    []interface{}{},
			"message": "Belum ada laporan pengiriman untuk tanggal ini.",
		})
	}

	// Ambil hanya item menu yang merupakan pengganti untuk siswa alergi
	substituteItems := make([]models.DeliveryMenuItem, 0)
	for _, item := range report.Items {
		if item.IsAllergySubstitute {
			substituteItems = append(substituteItems, item)
		}
	}

	// Kelompokkan per jenis alergi (dari AllergyNote)
	groupMap := make(map[string][]models.DeliveryMenuItem)
	for _, item := range substituteItems {
		groupMap[item.AllergyNote] = append(groupMap[item.AllergyNote], item)
	}

	// Untuk setiap kelompok, ambil jumlah siswa yang terdampak
	type MenuGroup struct {
		AllergyType  string                    `json:"allergy_type"`
		StudentCount int64                     `json:"student_count"`
		MenuItems    []models.DeliveryMenuItem `json:"menu_items"`
	}
	result := make([]MenuGroup, 0)
	for allergyType, items := range groupMap {
		var count int64
		config.DB.Model(&models.StudentAllergy{}).
			Where("school_id = ? AND allergy_type LIKE ?", userID, "%"+allergyType+"%").
			Count(&count)
		result = append(result, MenuGroup{
			AllergyType:  allergyType,
			StudentCount: count,
			MenuItems:    items,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"date":   targetDate.Format("2006-01-02"),
		"data":   result,
	})
}
