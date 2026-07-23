package sppg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
)

type mlPrediction struct {
	Calories            float64  `json:"calories"`
	Fat                 float64  `json:"fat"`
	Carbohydrate        float64  `json:"carbohydrate"`
	Protein             float64  `json:"protein"`
	StatusGizi          string   `json:"status_gizi"`
	DetectedIngredients []string `json:"detected_ingredients"`
	GeminiFailed        bool     `json:"gemini_failed"`
}

func callMLService(imageData []byte, filename string) (*mlPrediction, error) {
	mlURL := os.Getenv("ML_SERVICE_URL")
	if mlURL == "" {
		mlURL = "http://localhost:8000"
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(imageData); err != nil {
		return nil, err
	}
	writer.Close()

	// Timeout > worst-case ML (retry 503 Gemini + timeout HTTP 30s/attempt),
	// tapi tetap batasi agar request tak hang selamanya kalau ML macet.
	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Post(mlURL+"/predict", writer.FormDataContentType(), &buf)
	if err != nil {
		return nil, fmt.Errorf("ML service tidak dapat dijangkau: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ML service error %d: %s", resp.StatusCode, string(body))
	}

	var prediction mlPrediction
	if err := json.NewDecoder(resp.Body).Decode(&prediction); err != nil {
		return nil, fmt.Errorf("gagal parse response ML: %w", err)
	}
	return &prediction, nil
}

func ScanMenu(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	totalPortionsStr := c.FormValue("total_portions")
	totalPortions, err := strconv.Atoi(totalPortionsStr)
	if err != nil || totalPortions <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "total_portions harus diisi dan bernilai positif.",
		})
	}

	// Terima file gambar dari SPPG
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Field 'image' wajib diisi dengan file gambar.",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal membaca file gambar."})
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal memproses file gambar."})
	}

	// Panggil ML service, fallback ke dummy data jika tidak tersedia
	prediction, err := callMLService(imageData, fileHeader.Filename)
	if err != nil {
		prediction = &mlPrediction{
			Calories:            500,
			Fat:                 15,
			Carbohydrate:        70,
			Protein:             20,
			StatusGizi:          "Baik",
			DetectedIngredients: []string{},
			GeminiFailed:        true,
		}
	}

	detectedIngredients := prediction.DetectedIngredients

	// Ambil sekolah yang dilayani oleh SPPG ini
	var schoolProfiles []models.SchoolProfile
	config.DB.Where("sppg_user_id = ?", userID).Find(&schoolProfiles)

	schoolIDs := make([]uint, 0)
	for _, sp := range schoolProfiles {
		schoolIDs = append(schoolIDs, sp.UserID)
	}

	// Cocokkan bahan terdeteksi dengan data alergi siswa
	type IngredientStatus struct {
		Name       string `json:"name"`
		IsAllergen bool   `json:"is_allergen"`
	}
	ingredientStatuses := make([]IngredientStatus, 0)
	allergenNames := make([]string, 0)

	for _, ingredient := range detectedIngredients {
		isAllergen := false
		if len(schoolIDs) > 0 {
			var count int64
			config.DB.Model(&models.StudentAllergy{}).
				Where("school_id IN ? AND allergy_type ILIKE ?", schoolIDs, "%"+ingredient+"%").
				Count(&count)
			if count > 0 {
				isAllergen = true
				allergenNames = append(allergenNames, ingredient)
			}
		}
		ingredientStatuses = append(ingredientStatuses, IngredientStatus{
			Name:       ingredient,
			IsAllergen: isAllergen,
		})
	}

	// Bangun peringatan alergi per sekolah per alergen
	type AllergyWarning struct {
		AllergenType string `json:"allergen_type"`
		SchoolID     uint   `json:"school_id"`
		SchoolName   string `json:"school_name"`
		StudentCount int64  `json:"student_count"`
	}
	warnings := make([]AllergyWarning, 0)

	for _, allergen := range allergenNames {
		for _, sp := range schoolProfiles {
			var school models.User
			if config.DB.First(&school, sp.UserID).Error != nil {
				continue
			}
			var count int64
			config.DB.Model(&models.StudentAllergy{}).
				Where("school_id = ? AND allergy_type ILIKE ?", sp.UserID, "%"+allergen+"%").
				Count(&count)
			if count > 0 {
				warnings = append(warnings, AllergyWarning{
					AllergenType: strings.Title(allergen),
					SchoolID:     school.ID,
					SchoolName:   school.Name,
					StudentCount: count,
				})
			}
		}
	}

	kaloriPerPorsi := 0.0
	if totalPortions > 0 {
		kaloriPerPorsi = prediction.Calories / float64(totalPortions)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"menu_detail": fiber.Map{
				"kalori":      prediction.Calories,
				"protein":     prediction.Protein,
				"karbohidrat": prediction.Carbohydrate,
				"lemak":       prediction.Fat,
			},
			"kalori_per_porsi":  fmt.Sprintf("%.1f", kaloriPerPorsi),
			"total_porsi":       totalPortions,
			"status_gizi":       prediction.StatusGizi,
			"bahan_terdeteksi":  ingredientStatuses,
			"peringatan_alergi": warnings,
			"gemini_failed":     prediction.GeminiFailed,
		},
	})
}

func SubmitMenuReport(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	type submitBody struct {
		MenuName       string  `json:"menu_name"`
		Kalori         float64 `json:"kalori"`
		Protein        float64 `json:"protein"`
		Karbohidrat    float64 `json:"karbohidrat"`
		Lemak          float64 `json:"lemak"`
		TotalPorsi     int     `json:"total_porsi"`
		KaloriPerPorsi float64 `json:"kalori_per_porsi"`
		StatusGizi     string  `json:"status_gizi"`
	}

	var body submitBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format body tidak valid."})
	}

	if body.TotalPorsi <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "total_porsi harus bernilai positif."})
	}

	report := models.MenuReport{
		SppgID:         userID,
		MenuName:       body.MenuName,
		Kalori:         body.Kalori,
		Protein:        body.Protein,
		Karbohidrat:    body.Karbohidrat,
		Lemak:          body.Lemak,
		TotalPorsi:     body.TotalPorsi,
		KaloriPerPorsi: body.KaloriPerPorsi,
		StatusGizi:     body.StatusGizi,
		ReportDate:     time.Now(),
	}

	if err := config.DB.Create(&report).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan laporan menu."})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Laporan menu berhasil dikirim.",
		"data":    report,
	})
}
