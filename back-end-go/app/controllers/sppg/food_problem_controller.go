package sppg

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var validPriorities = map[string]bool{
	"Rendah": true,
	"Sedang": true,
	"Tinggi": true,
}

// CreateFoodProblem dipakai SPPG untuk melaporkan masalah pada makanan
// yang ditemukan secara internal (sebelum atau saat distribusi).
func CreateFoodProblem(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	schoolName := strings.TrimSpace(c.FormValue("namaSekolah"))
	foodName := strings.TrimSpace(c.FormValue("namaMakanan"))
	problemType := strings.TrimSpace(c.FormValue("jenisMasalah"))
	priority := strings.TrimSpace(c.FormValue("prioritas"))
	description := strings.TrimSpace(c.FormValue("deskripsiMasalah"))
	portionsStr := strings.TrimSpace(c.FormValue("porsiTerdampak"))

	if schoolName == "" || foodName == "" || problemType == "" || description == "" || portionsStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Lengkapi semua data wajib (nama sekolah, nama makanan, jenis masalah, porsi, deskripsi).",
		})
	}

	portions, err := strconv.Atoi(portionsStr)
	if err != nil || portions <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "porsiTerdampak harus berupa angka > 0."})
	}

	if priority == "" {
		priority = "Sedang"
	}
	if !validPriorities[priority] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "prioritas tidak valid. Pilihan: Rendah, Sedang, Tinggi."})
	}

	photoPath := ""
	if photo, err := c.FormFile("fotoMasalah"); err == nil {
		uploadDir := "./uploads/sppg-problems"
		os.MkdirAll(uploadDir, os.ModePerm)
		filename := fmt.Sprintf("sppg_problem_%d_%d_%s", userID, time.Now().UnixNano(), filepath.Base(photo.Filename))
		path := filepath.Join(uploadDir, filename)
		if err := c.SaveFile(photo, path); err == nil {
			photoPath = path
		}
	}

	problem := models.SppgFoodProblem{
		SppgID:           userID,
		SchoolName:       schoolName,
		FoodName:         foodName,
		ProblemType:      problemType,
		AffectedPortions: portions,
		Priority:         priority,
		Description:      description,
		PhotoPath:        photoPath,
		Status:           "baru",
	}

	if err := config.DB.Create(&problem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan laporan."})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Laporan masalah makanan berhasil dikirim.",
		"data":    problem,
	})
}

// GetSchoolFoodReports mengembalikan daftar laporan makanan bermasalah yang
// dikirim oleh sekolah ke SPPG ini (sumber: tabel FoodReport).
func GetSchoolFoodReports(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	month := c.Query("month")
	year := c.Query("year")

	query := config.DB.Where("sppg_id = ?", userID).Order("report_date desc")
	if month != "" {
		query = query.Where("EXTRACT(MONTH FROM report_date) = ?", month)
	}
	if year != "" {
		query = query.Where("EXTRACT(YEAR FROM report_date) = ?", year)
	}

	var reports []models.FoodReport
	if err := query.Find(&reports).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil data laporan."})
	}

	type Row struct {
		ID          uint   `json:"id"`
		ReportDate  string `json:"report_date"`
		Reporter    string `json:"reporter"`
		IssueType   string `json:"issue_type"`
		Description string `json:"description"`
		PhotoURL    string `json:"photo_url"`
		Status      string `json:"status"`
	}
	out := make([]Row, 0, len(reports))
	for _, r := range reports {
		var school models.User
		config.DB.Select("id, name").First(&school, r.SchoolID)
		out = append(out, Row{
			ID:          r.ID,
			ReportDate:  r.ReportDate.Format("2006-01-02"),
			Reporter:    school.Name,
			IssueType:   r.IssueType,
			Description: r.Description,
			PhotoURL:    r.PhotoURL,
			Status:      r.Status,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   out,
	})
}
