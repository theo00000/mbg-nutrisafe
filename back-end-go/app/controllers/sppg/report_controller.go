package sppg

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DeliveryMenuItemInput struct {
	MenuName            string `json:"menu_name"`
	Category            string `json:"category"`
	Portions            int    `json:"portions"`
	IsAllergySubstitute bool   `json:"is_allergy_substitute"`
	AllergyNote         string `json:"allergy_note"`
}

func GetDeliveryReports(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	month := c.Query("month")
	year := c.Query("year")

	query := config.DB.Where("sppg_id = ?", userID).Order("delivery_date desc")
	if month != "" && year != "" {
		query = query.Where("MONTH(delivery_date) = ? AND YEAR(delivery_date) = ?", month, year)
	} else if year != "" {
		query = query.Where("YEAR(delivery_date) = ?", year)
	}

	var reports []models.DeliveryReport
	if err := query.Find(&reports).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil data laporan."})
	}

	// Kelompokkan per tanggal
	type SchoolEntry struct {
		ReportID      uint   `json:"report_id"`
		SchoolID      uint   `json:"school_id"`
		SchoolName    string `json:"school_name"`
		Address       string `json:"address"`
		TotalPortions int    `json:"total_portions"`
		Status        string `json:"status"`
	}
	dateMap := make(map[string][]SchoolEntry)
	for _, r := range reports {
		var school models.User
		config.DB.First(&school, r.SchoolID)
		var profile models.SchoolProfile
		config.DB.Where("user_id = ?", r.SchoolID).First(&profile)

		dateKey := r.DeliveryDate.Format("2006-01-02")
		dateMap[dateKey] = append(dateMap[dateKey], SchoolEntry{
			ReportID:      r.ID,
			SchoolID:      r.SchoolID,
			SchoolName:    school.Name,
			Address:       profile.Address,
			TotalPortions: r.TotalPortions,
			Status:        r.Status,
		})
	}

	result := make([]fiber.Map, 0)
	for date, schools := range dateMap {
		result = append(result, fiber.Map{
			"date":    date,
			"schools": schools,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func GetDeliveryReportDetail(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	id := c.Params("id")
	var report models.DeliveryReport
	if err := config.DB.Preload("Items").Where("id = ? AND sppg_id = ?", id, userID).First(&report).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Laporan tidak ditemukan."})
	}

	var school models.User
	config.DB.First(&school, report.SchoolID)
	var profile models.SchoolProfile
	config.DB.Where("user_id = ?", report.SchoolID).First(&profile)

	type AllergyCount struct {
		AllergyType string `json:"allergy_type"`
		Count       int    `json:"count"`
	}
	var allergyStats []AllergyCount
	config.DB.Model(&models.StudentAllergy{}).
		Select("allergy_type, count(*) as count").
		Where("school_id = ?", report.SchoolID).
		Group("allergy_type").
		Scan(&allergyStats)

	var totalAllergyStudents int64
	config.DB.Model(&models.StudentAllergy{}).Where("school_id = ?", report.SchoolID).Count(&totalAllergyStudents)

	distributionPct := 0
	if report.TotalPortions > 0 {
		distributionPct = (report.DistributedPortions * 100) / report.TotalPortions
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id": report.ID,
			"school": fiber.Map{
				"id":      school.ID,
				"name":    school.Name,
				"address": profile.Address,
			},
			"delivery_date":        report.DeliveryDate.Format("2006-01-02"),
			"status":               report.Status,
			"total_portions":       report.TotalPortions,
			"distributed_portions": report.DistributedPortions,
			"distribution_pct":     distributionPct,
			"allergy_students": fiber.Map{
				"total":     totalAllergyStudents,
				"breakdown": allergyStats,
			},
			"menu_items": report.Items,
			"photos": fiber.Map{
				"photo1": report.Photo1Path,
				"photo2": report.Photo2Path,
			},
		},
	})
}

func CreateDeliveryReport(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	schoolID, err := strconv.Atoi(c.FormValue("school_id"))
	if err != nil || schoolID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "school_id tidak valid."})
	}

	deliveryDate, err := time.Parse("2006-01-02", c.FormValue("delivery_date"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format tanggal tidak valid. Gunakan YYYY-MM-DD."})
	}

	totalPortions, err := strconv.Atoi(c.FormValue("total_portions"))
	if err != nil || totalPortions <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "total_portions tidak valid."})
	}

	distributedPortions := totalPortions
	if dp := c.FormValue("distributed_portions"); dp != "" {
		if v, err := strconv.Atoi(dp); err == nil {
			distributedPortions = v
		}
	}

	var menuItems []DeliveryMenuItemInput
	if raw := c.FormValue("menu_items"); raw != "" {
		if err := json.Unmarshal([]byte(raw), &menuItems); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format menu_items tidak valid."})
		}
	}

	uploadDir := "./uploads/delivery"
	os.MkdirAll(uploadDir, os.ModePerm)

	photo1Path, photo2Path := "", ""
	if f, err := c.FormFile("photo1"); err == nil {
		name := fmt.Sprintf("delivery_%d_%d_%s", userID, time.Now().UnixNano(), filepath.Base(f.Filename))
		path := filepath.Join(uploadDir, name)
		if err := c.SaveFile(f, path); err == nil {
			photo1Path = path
		}
	}
	if f, err := c.FormFile("photo2"); err == nil {
		name := fmt.Sprintf("delivery_%d_%d_%s", userID, time.Now().UnixNano(), filepath.Base(f.Filename))
		path := filepath.Join(uploadDir, name)
		if err := c.SaveFile(f, path); err == nil {
			photo2Path = path
		}
	}

	report := models.DeliveryReport{
		SppgID:              userID,
		SchoolID:            uint(schoolID),
		DeliveryDate:        deliveryDate,
		TotalPortions:       totalPortions,
		DistributedPortions: distributedPortions,
		Status:              "pending",
		Photo1Path:          photo1Path,
		Photo2Path:          photo2Path,
	}

	tx := config.DB.Begin()
	if err := tx.Create(&report).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan laporan."})
	}
	for _, item := range menuItems {
		mi := models.DeliveryMenuItem{
			DeliveryReportID:    report.ID,
			MenuName:            item.MenuName,
			Category:            item.Category,
			Portions:            item.Portions,
			IsAllergySubstitute: item.IsAllergySubstitute,
			AllergyNote:         item.AllergyNote,
		}
		if err := tx.Create(&mi).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan item menu."})
		}
	}
	tx.Commit()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Laporan pengiriman berhasil dikirim.",
		"data": fiber.Map{
			"id":             report.ID,
			"school_id":      report.SchoolID,
			"delivery_date":  deliveryDate.Format("2006-01-02"),
			"total_portions": report.TotalPortions,
			"status":         report.Status,
		},
	})
}
