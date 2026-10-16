package school

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/app/services"
	"back-end/config"
	"back-end/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

var validIssueTypes = map[string]bool{
	"kekurangan_makanan": true,
	"makanan_basi":       true,
	"masalah_alergen":    true,
	"terlambat":          true,
	"lainnya":            true,
}

// GetAssignedsppg mengembalikan info sppg yang ditugaskan ke sekolah ini (untuk halaman pelaporan)
func GetAssignedsppg(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var profile models.SchoolProfile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Profil sekolah tidak ditemukan."})
	}

	if profile.SppgUserID == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"data":    nil,
			"message": "Sekolah belum memiliki mitra sppg",
		})
	}

	var sppg models.User
	if err := config.DB.First(&sppg, *profile.SppgUserID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Data sppg tidak ditemukan."})
	}

	var sppgProfile models.SppgProfile
	config.DB.Where("user_id = ?", sppg.ID).First(&sppgProfile)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":       sppg.ID,
			"name":     sppg.Name,
			"address":  sppgProfile.Address,
			"capacity": sppgProfile.Capacity,
			"status":   sppgProfile.Status,
		},
	})
}

func CreateFoodReport(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var profile models.SchoolProfile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Profil sekolah tidak ditemukan"})
	}
	if profile.SppgUserID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Sekolah belum memiliki mitra sppg"})
	}

	issueType := c.FormValue("issue_type")
	description := c.FormValue("description")
	reportDateStr := c.FormValue("report_date")

	if issueType == "" || description == "" || reportDateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "issue_type, description, dan report_date wajib diisi"})
	}
	if !validIssueTypes[issueType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "issue_type tidak valid. Pilihan: kekurangan_makanan, makanan_basi, masalah_alergen, terlambat, lainnya.",
		})
	}

	reportDate, err := time.Parse("2006-01-02", reportDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format tanggal tidak valid. Gunakan YYYY-MM-DD."})
	}

	photoURL := ""
	if photo, err := c.FormFile("photo"); err == nil {
		uploadDir := "./uploads/reports"
		os.MkdirAll(uploadDir, os.ModePerm)
		filename := fmt.Sprintf("report_%d_%d_%s", userID, time.Now().UnixNano(), filepath.Base(photo.Filename))
		path := filepath.Join(uploadDir, filename)
		if err := c.SaveFile(photo, path); err == nil {
			photoURL = path
		}
	}

	report := models.FoodReport{
		SchoolID:    userID,
		SppgID:      *profile.SppgUserID,
		ReporterID:  userID,
		IssueType:   issueType,
		Description: description,
		PhotoURL:    photoURL,
		ReportDate:  reportDate,
		Status:      "belum_ditinjau",
	}

	if err := config.DB.Create(&report).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan laporan."})
	}

	notifyFoodReportCreated(report)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Laporan berhasil dikirim.",
		"data":    services.BuildFoodReportDetail(report),
	})
}

// notifyFoodReportCreated memberi tahu pengirim bahwa laporannya masuk, dan
// meneruskan laporan itu ke SPPG terkait serta admin BGN.
func notifyFoodReportCreated(report models.FoodReport) {
	var sppgUser models.User
	config.DB.Select("id, name").First(&sppgUser, report.SppgID)
	var schoolUser models.User
	config.DB.Select("id, name").First(&schoolUser, report.SchoolID)

	tanggal := utils.TanggalID(report.ReportDate)

	services.Notify(
		report.ReporterID,
		models.NotifLaporanTerkirim,
		"Laporan Terkirim!",
		fmt.Sprintf("Laporan makanan kepada %s pada tanggal %s berhasil terkirim!", sppgUser.Name, tanggal),
		"food_report", report.ID,
	)

	services.Notify(
		report.SppgID,
		models.NotifLaporanMasuk,
		"Laporan Masalah Baru",
		fmt.Sprintf("%s melaporkan %s pada tanggal %s.",
			schoolUser.Name, services.IssueTypeLabel(report.IssueType), tanggal),
		"food_report", report.ID,
	)

	services.NotifyAdmins(
		models.NotifLaporanMasuk,
		"Laporan Masalah Baru",
		fmt.Sprintf("%s melaporkan %s terhadap %s pada tanggal %s.",
			schoolUser.Name, services.IssueTypeLabel(report.IssueType), sppgUser.Name, tanggal),
		"food_report", report.ID,
	)
}

// GetFoodReportDetail mengembalikan satu laporan milik sekolah ini.
// Penyaringan school_id memastikan sekolah lain tidak bisa membuka laporannya.
func GetFoodReportDetail(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var report models.FoodReport
	if err := config.DB.Where("id = ? AND school_id = ?", c.Params("id"), userID).First(&report).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Laporan tidak ditemukan."})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   services.BuildFoodReportDetail(report),
	})
}

func GetFoodReports(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var reports []models.FoodReport
	if err := config.DB.Where("school_id = ?", userID).Order("created_at desc").Find(&reports).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil data laporan."})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   services.BuildFoodReportDetails(reports),
	})
}
