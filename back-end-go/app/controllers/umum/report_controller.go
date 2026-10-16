package umum

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/app/services"
	"back-end/config"
	"back-end/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetSiswaSchool mengembalikan info sekolah yang terhubung ke akun siswa yang sedang login.
func GetSiswaSchool(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var student models.Student
	if err := config.DB.Where("user_id = ?", userID).First(&student).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Data siswa tidak ditemukan"})
	}

	var schoolProfile models.SchoolProfile
	if err := config.DB.Where("user_id = ?", student.SchoolID).First(&schoolProfile).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Profil sekolah tidak ditemukan"})
	}

	var schoolUser models.User
	if err := config.DB.First(&schoolUser, student.SchoolID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Akun sekolah tidak ditemukan"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"school_id":   student.SchoolID,
			"school_name": schoolUser.Name,
			"address":     schoolProfile.Address,
		},
	})
}

func GetSchoolsList(c *fiber.Ctx) error {
	var schools []models.SchoolProfile
	if err := config.DB.Find(&schools).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil data sekolah"})
	}

	schoolIDs := make([]uint, len(schools))
	for i, s := range schools {
		schoolIDs[i] = s.UserID
	}

	var users []models.User
	config.DB.Where("id IN ?", schoolIDs).Find(&users)
	userMap := make(map[uint]models.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	result := make([]fiber.Map, 0, len(schools))
	for _, s := range schools {
		u := userMap[s.UserID]
		result = append(result, fiber.Map{
			"school_id":   s.UserID,
			"school_name": u.Name,
			"address":     s.Address,
			"Grade":       s.Grade,
			"npsn":        s.NPSN,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func CreatePublicReport(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	schoolIDRaw := c.FormValue("school_id")
	issueType := c.FormValue("issue_type")
	reportDateStr := c.FormValue("report_date")
	description := c.FormValue("description")

	if schoolIDRaw == "" || issueType == "" || reportDateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "school_id, issue_type, dan report_date wajib diisi"})
	}

	validIssues := map[string]bool{
		"kekurangan_makanan": true,
		"makanan_basi":       true,
		"masalah_alergen":    true,
		"terlambat":          true,
		"lainnya":            true,
	}
	if !validIssues[issueType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Jenis masalah tidak valid"})
	}

	reportDate, err := time.Parse("2006-01-02", reportDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format tanggal tidak valid. Gunakan YYYY-MM-DD."})
	}

	var schoolIDUint uint
	fmt.Sscanf(schoolIDRaw, "%d", &schoolIDUint)

	var schoolProfile models.SchoolProfile
	if err := config.DB.Where("user_id = ?", schoolIDUint).First(&schoolProfile).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Sekolah tidak ditemukan"})
	}

	var sppgID uint
	if schoolProfile.SppgUserID != nil {
		sppgID = *schoolProfile.SppgUserID
	}

	var photoPath string
	if file, err := c.FormFile("photo"); err == nil {
		// Tanpa MkdirAll, SaveFile gagal saat direktori belum ada dan foto
		// hilang tanpa jejak — error-nya dulu ditelan diam-diam.
		uploadDir := "./uploads/reports"
		if mkErr := os.MkdirAll(uploadDir, os.ModePerm); mkErr != nil {
			log.Printf("[upload] gagal menyiapkan %s: %v", uploadDir, mkErr)
		}
		filename := fmt.Sprintf("public_report_%d_%d%s", userID, time.Now().UnixNano(), filepath.Ext(file.Filename))
		savePath := filepath.Join(uploadDir, filename)
		if saveErr := c.SaveFile(file, savePath); saveErr != nil {
			log.Printf("[upload] gagal menyimpan foto laporan user=%d: %v", userID, saveErr)
		} else {
			photoPath = savePath
		}
	}

	report := models.FoodReport{
		SchoolID:    schoolIDUint,
		SppgID:      sppgID,
		ReporterID:  userID,
		IssueType:   issueType,
		Description: description,
		PhotoURL:    photoPath,
		ReportDate:  reportDate,
		Status:      "belum_ditinjau",
	}

	if err := config.DB.Create(&report).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan laporan"})
	}

	notifyPublicReportCreated(report)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Laporan berhasil dikirim",
		"data":    services.BuildFoodReportDetail(report),
	})
}

// notifyPublicReportCreated mengabari pelapor, lalu meneruskan ke sekolah,
// SPPG (kalau sekolah sudah punya mitra), dan admin.
func notifyPublicReportCreated(report models.FoodReport) {
	var reporter, schoolUser, sppgUser models.User
	config.DB.Select("id, name").First(&reporter, report.ReporterID)
	config.DB.Select("id, name").First(&schoolUser, report.SchoolID)
	if report.SppgID != 0 {
		config.DB.Select("id, name").First(&sppgUser, report.SppgID)
	}

	tanggal := utils.TanggalID(report.ReportDate)
	issue := services.IssueTypeLabel(report.IssueType)

	tujuan := schoolUser.Name
	if sppgUser.Name != "" {
		tujuan = sppgUser.Name
	}
	services.Notify(
		report.ReporterID,
		models.NotifLaporanTerkirim,
		"Laporan Terkirim!",
		fmt.Sprintf("Laporan makanan kepada %s pada tanggal %s berhasil terkirim!", tujuan, tanggal),
		"food_report", report.ID,
	)

	services.Notify(
		report.SchoolID,
		models.NotifLaporanMasuk,
		"Laporan Masyarakat Baru",
		fmt.Sprintf("%s melaporkan %s pada tanggal %s.", reporter.Name, issue, tanggal),
		"food_report", report.ID,
	)

	// SppgID 0 berarti sekolah belum punya mitra — Notify mengabaikannya.
	services.Notify(
		report.SppgID,
		models.NotifLaporanMasuk,
		"Laporan Masyarakat Baru",
		fmt.Sprintf("%s melaporkan %s di %s pada tanggal %s.", reporter.Name, issue, schoolUser.Name, tanggal),
		"food_report", report.ID,
	)

	services.NotifyAdmins(
		models.NotifLaporanMasuk,
		"Laporan Masyarakat Baru",
		fmt.Sprintf("%s melaporkan %s di %s pada tanggal %s.", reporter.Name, issue, schoolUser.Name, tanggal),
		"food_report", report.ID,
	)
}

// GetPublicReports mengembalikan riwayat laporan yang dikirim akun ini.
// Sebelumnya laporan umum tidak bisa dilihat lagi setelah terkirim.
func GetPublicReports(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var reports []models.FoodReport
	if err := config.DB.Where("reporter_id = ?", userID).Order("created_at desc").Find(&reports).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil data laporan"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   services.BuildFoodReportDetails(reports),
	})
}

// GetPublicReportDetail membuka satu laporan milik pelapor ini saja.
func GetPublicReportDetail(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var report models.FoodReport
	if err := config.DB.Where("id = ? AND reporter_id = ?", c.Params("id"), userID).First(&report).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Laporan tidak ditemukan"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   services.BuildFoodReportDetail(report),
	})
}
