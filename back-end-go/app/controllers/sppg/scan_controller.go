package sppg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/app/services"
	"back-end/config"
	"back-end/utils"

	"github.com/gofiber/fiber/v2"
)

const scanUploadDir = "./uploads/menu-scan"

type mlPrediction struct {
	Calories            float64  `json:"calories"`
	Fat                 float64  `json:"fat"`
	Carbohydrate        float64  `json:"carbohydrate"`
	Protein             float64  `json:"protein"`
	StatusGizi          string   `json:"status_gizi"`
	DetectedIngredients []string `json:"detected_ingredients"`
	GeminiFailed        bool     `json:"gemini_failed"`
}

type mlRecommendation struct {
	Allergen           string `json:"allergen"`
	OriginalIngredient string `json:"original_ingredient"`
	Substitute         string `json:"substitute"`
	Category           string `json:"category"`
	Reason             string `json:"reason"`
}

type mlRecommendResponse struct {
	Recommendations []mlRecommendation `json:"recommendations"`
	GeminiFailed    bool               `json:"gemini_failed"`
}

func mlBaseURL() string {
	if u := os.Getenv("ML_SERVICE_URL"); u != "" {
		return u
	}
	return "http://localhost:8000"
}

func callMLService(imageData []byte, filename string) (*mlPrediction, error) {
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
	resp, err := client.Post(mlBaseURL()+"/predict", writer.FormDataContentType(), &buf)
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

// recommendEnabled menandai apakah endpoint rekomendasi di ML service sudah
// tersedia.
//
// Model rekomendasi masih dikerjakan dan belum diputuskan apakah memakai
// Gemini atau model ML murni — bentuk responsnya bisa berubah. Sampai itu
// final, fitur ini mati secara default supaya tidak ada permintaan ke
// endpoint yang belum ada. Nyalakan dengan ML_RECOMMEND_ENABLED=true, dan
// cocokkan dulu mlRecommendResponse dengan kontrak yang benar-benar dipakai.
func recommendEnabled() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("ML_RECOMMEND_ENABLED")), "true")
}

// callMLRecommend meminta bahan pengganti untuk alergen yang terdeteksi.
// Kegagalan dikembalikan apa adanya; pemanggil memutuskan apakah scan tetap
// diteruskan tanpa rekomendasi (lihat ScanMenu).
func callMLRecommend(menuName string, allergens, ingredients []string) (*mlRecommendResponse, error) {
	if !recommendEnabled() {
		return nil, fmt.Errorf("endpoint rekomendasi ML belum tersedia")
	}

	payload, err := json.Marshal(map[string]interface{}{
		"menu_name":            menuName,
		"allergens":            allergens,
		"detected_ingredients": ingredients,
	})
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Post(mlBaseURL()+"/recommend", "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("ML service tidak dapat dijangkau: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ML service error %d: %s", resp.StatusCode, string(body))
	}

	var out mlRecommendResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("gagal parse response rekomendasi: %w", err)
	}
	return &out, nil
}

// titleCase menggantikan strings.Title yang sudah deprecated. Hanya huruf
// pertama tiap kata yang dibesarkan, sisanya dibiarkan apa adanya.
func titleCase(s string) string {
	prevIsSpace := true
	return strings.Map(func(r rune) rune {
		if prevIsSpace && unicode.IsLetter(r) {
			prevIsSpace = false
			return unicode.ToUpper(r)
		}
		prevIsSpace = unicode.IsSpace(r)
		return r
	}, s)
}

// sppgSchoolIDs mengembalikan user_id sekolah yang dilayani SPPG ini.
func sppgSchoolIDs(sppgID uint) ([]models.SchoolProfile, []uint) {
	var profiles []models.SchoolProfile
	config.DB.Where("sppg_user_id = ?", sppgID).Find(&profiles)

	ids := make([]uint, 0, len(profiles))
	for _, p := range profiles {
		ids = append(ids, p.UserID)
	}
	return profiles, ids
}

type ingredientStatus struct {
	Name       string `json:"name"`
	IsAllergen bool   `json:"is_allergen"`
}

// Kata yang tidak membawa makna saat mencocokkan alergi dengan bahan.
// Sekolah menulis jenis alergi dengan bebas ("Alergi Kacang", "Alergi
// terhadap Susu"), jadi kata-kata ini harus dibuang lebih dulu.
var allergyStopwords = map[string]bool{
	"alergi": true, "dan": true, "atau": true, "terhadap": true,
	"pada": true, "makanan": true, "jenis": true,
}

// tokenize memecah teks jadi kata bermakna dalam huruf kecil.
// Kata di bawah 3 huruf dibuang supaya tidak memicu kecocokan palsu.
func tokenize(s string) []string {
	fields := strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		if len([]rune(f)) < 3 || allergyStopwords[f] {
			continue
		}
		out = append(out, f)
	}
	return out
}

// tokensMatch benar kalau ada satu kata yang sama-sama muncul di kedua sisi.
//
// Pencocokan substring pada kalimat utuh tidak cukup: "Alergi Kacang" tidak
// mengandung "Kacang Tanah" dan sebaliknya, padahal keduanya jelas merujuk
// kacang yang sama. Per-kata juga menangkap bentuk berimbuhan seperti
// "kacang-kacangan".
func tokensMatch(a, b []string) bool {
	for _, x := range a {
		for _, y := range b {
			if x == y || strings.Contains(x, y) || strings.Contains(y, x) {
				return true
			}
		}
	}
	return false
}

// classifyIngredients menandai bahan mana yang cocok dengan data alergi siswa
// di sekolah binaan SPPG. Dipakai saat scan maupun saat submit, supaya status
// alergen selalu dihitung server dan tidak bergantung pada kiriman klien.
//
// Mengembalikan status tiap bahan, daftar bahan yang memicu alergi, dan
// pemetaan bahan -> jenis alergi persis yang cocok. Pemetaan itu dipakai
// buildAllergyWarnings agar bisa menghitung siswa dengan pencarian tepat,
// bukan ILIKE yang punya kelemahan sama dengan pencocokan kalimat utuh.
func classifyIngredients(ingredients []string, schoolIDs []uint) ([]ingredientStatus, []string, map[string][]string) {
	statuses := make([]ingredientStatus, 0, len(ingredients))
	allergens := make([]string, 0)
	matched := make(map[string][]string)

	if len(ingredients) == 0 || len(schoolIDs) == 0 {
		for _, name := range ingredients {
			statuses = append(statuses, ingredientStatus{Name: name, IsAllergen: false})
		}
		return statuses, allergens, matched
	}

	var allergyTypes []string
	config.DB.Model(&models.StudentAllergy{}).
		Where("school_id IN ?", schoolIDs).
		Distinct("allergy_type").
		Pluck("allergy_type", &allergyTypes)

	type allergyEntry struct {
		raw    string
		tokens []string
	}
	entries := make([]allergyEntry, 0, len(allergyTypes))
	for _, t := range allergyTypes {
		if tokens := tokenize(t); len(tokens) > 0 {
			entries = append(entries, allergyEntry{raw: t, tokens: tokens})
		}
	}

	for _, name := range ingredients {
		nameTokens := tokenize(name)
		hits := make([]string, 0)
		for _, e := range entries {
			if tokensMatch(nameTokens, e.tokens) {
				hits = append(hits, e.raw)
			}
		}
		if len(hits) > 0 {
			allergens = append(allergens, name)
			matched[name] = hits
		}
		statuses = append(statuses, ingredientStatus{Name: name, IsAllergen: len(hits) > 0})
	}

	return statuses, allergens, matched
}

type allergyWarning struct {
	AllergenType string `json:"allergen_type"`
	SchoolID     uint   `json:"school_id"`
	SchoolName   string `json:"school_name"`
	StudentCount int64  `json:"student_count"`
}

// buildAllergyWarnings merinci sekolah mana dan berapa siswa yang terdampak
// tiap alergen, untuk kartu "Peringatan Alergi".
//
// matched memetakan nama bahan ke jenis alergi persis yang sudah dicocokkan
// classifyIngredients, jadi penghitungan siswa memakai pencarian tepat.
// Memakai ILIKE "%bahan%" di sini akan salah dengan cara yang sama seperti
// pencocokan kalimat utuh: "Alergi Kacang" tidak cocok "%Kacang Tanah%".
func buildAllergyWarnings(allergens []string, matched map[string][]string, profiles []models.SchoolProfile) []allergyWarning {
	warnings := make([]allergyWarning, 0)
	if len(allergens) == 0 || len(profiles) == 0 {
		return warnings
	}

	schoolIDs := make([]uint, 0, len(profiles))
	for _, p := range profiles {
		schoolIDs = append(schoolIDs, p.UserID)
	}

	// Satu query untuk semua nama sekolah, bukan satu per baris peringatan.
	var users []models.User
	config.DB.Select("id, name").Where("id IN ?", schoolIDs).Find(&users)
	nameByID := make(map[uint]string, len(users))
	for _, u := range users {
		nameByID[u.ID] = u.Name
	}

	for _, allergen := range allergens {
		types := matched[allergen]
		if len(types) == 0 {
			continue
		}

		type row struct {
			SchoolID uint
			Total    int64
		}
		var rows []row
		config.DB.Model(&models.StudentAllergy{}).
			Select("school_id, count(*) as total").
			Where("school_id IN ? AND allergy_type IN ?", schoolIDs, types).
			Group("school_id").
			Scan(&rows)

		for _, r := range rows {
			if r.Total == 0 {
				continue
			}
			warnings = append(warnings, allergyWarning{
				AllergenType: titleCase(allergen),
				SchoolID:     r.SchoolID,
				SchoolName:   nameByID[r.SchoolID],
				StudentCount: r.Total,
			})
		}
	}

	return warnings
}

func ScanMenu(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	totalPortions, err := strconv.Atoi(c.FormValue("total_portions"))
	if err != nil || totalPortions <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "total_portions harus diisi dan bernilai positif.",
		})
	}

	menuName := strings.TrimSpace(c.FormValue("menu_name"))

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

	// Simpan foto sekarang supaya saat submit laporan klien tidak perlu
	// mengunggah ulang; yang dikirim balik cuma nama file (lihat resolveScanPhoto).
	photoRef := ""
	os.MkdirAll(scanUploadDir, os.ModePerm)
	photoName := fmt.Sprintf("menu_%d_%d%s", userID, time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	if err := c.SaveFile(fileHeader, filepath.Join(scanUploadDir, photoName)); err == nil {
		photoRef = photoName
	}

	// Panggil ML service, fallback ke dummy data jika tidak tersedia.
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

	profiles, schoolIDs := sppgSchoolIDs(userID)
	ingredientStatuses, allergenNames, matchedTypes := classifyIngredients(prediction.DetectedIngredients, schoolIDs)
	warnings := buildAllergyWarnings(allergenNames, matchedTypes, profiles)

	// Kalau ada alergen, langsung ambil usulan penggantinya supaya layar hasil
	// bisa menampilkan kartu "Rekomendasi Menu Siswa Alergi" tanpa round-trip
	// kedua. Rekomendasi gagal tidak membatalkan hasil scan.
	//
	// Selama model rekomendasi belum siap, blok ini dilewati sepenuhnya —
	// tidak ada permintaan ke ML. Klien membaca "rekomendasi_aktif" untuk
	// membedakan "fitur belum ada" dari "rekomendasi gagal diambil".
	recommendations := make([]mlRecommendation, 0)
	recommendationFailed := false
	if len(allergenNames) > 0 && recommendEnabled() {
		rec, recErr := callMLRecommend(menuName, allergenNames, prediction.DetectedIngredients)
		if recErr != nil || rec == nil {
			recommendationFailed = true
		} else {
			recommendations = rec.Recommendations
			recommendationFailed = rec.GeminiFailed
		}
	}

	kaloriPerPorsi := prediction.Calories / float64(totalPortions)

	allergenStatus := models.AllergenStatusBebas
	if len(allergenNames) > 0 {
		allergenStatus = models.AllergenStatusMengandung
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"menu_name": menuName,
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
			"allergen_status":   allergenStatus,
			"rekomendasi":       recommendations,
			"rekomendasi_gagal": recommendationFailed,
			"rekomendasi_aktif": recommendEnabled(),
			"photo_ref":         photoRef,
			"photo_url":         scanPhotoURL(photoRef),
			"gemini_failed":     prediction.GeminiFailed,
		},
	})
}

// RecommendMenu memungkinkan klien meminta ulang rekomendasi tanpa mengunggah
// gambar lagi, mis. setelah pengguna mengoreksi nama menu.
func RecommendMenu(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	// Model rekomendasi masih dikerjakan. Jawab tegas supaya klien tidak
	// menampilkan kartu rekomendasi kosong seolah-olah tidak ada saran.
	if !recommendEnabled() {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":            "error",
			"message":           "Fitur rekomendasi menu belum tersedia.",
			"rekomendasi_aktif": false,
		})
	}

	var body struct {
		MenuName    string   `json:"menu_name"`
		Ingredients []string `json:"ingredients"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format body tidak valid."})
	}
	if len(body.Ingredients) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ingredients wajib diisi."})
	}

	_, schoolIDs := sppgSchoolIDs(userID)
	_, allergenNames, _ := classifyIngredients(body.Ingredients, schoolIDs)

	if len(allergenNames) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Tidak ada bahan yang memicu alergi siswa.",
			"data": fiber.Map{
				"allergen_status": models.AllergenStatusBebas,
				"rekomendasi":     []mlRecommendation{},
			},
		})
	}

	rec, err := callMLRecommend(body.MenuName, allergenNames, body.Ingredients)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  "error",
			"message": "Layanan rekomendasi sedang tidak tersedia. Coba lagi beberapa saat.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"allergen_status":   models.AllergenStatusMengandung,
			"alergen":           allergenNames,
			"rekomendasi":       rec.Recommendations,
			"rekomendasi_gagal": rec.GeminiFailed,
		},
	})
}

// scanPhotoURL mengembalikan URL foto hasil scan. Ref kosong (simpan gagal)
// tetap kosong, bukan URL direktori.
func scanPhotoURL(ref string) string {
	if ref == "" {
		return ""
	}
	return utils.PublicURL(filepath.Join(scanUploadDir, ref))
}

// resolveScanPhoto mengubah photo_ref dari klien menjadi path yang aman.
// Hanya nama file polos yang diterima — tanpa ini klien bisa menitipkan
// path sembarang (mis. "../../.env") untuk disimpan ke basis data.
func resolveScanPhoto(ref string) string {
	ref = strings.TrimSpace(ref)
	if ref == "" || ref != filepath.Base(ref) || strings.ContainsAny(ref, `/\`) {
		return ""
	}
	path := filepath.Join(scanUploadDir, ref)
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return path
}

// GetTodayMenuReport dipakai layar scan untuk tahu apakah laporan hari ini
// sudah dikirim ("Laporan Hari Ini Telah di Kirim!") sebelum menampilkan form.
func GetTodayMenuReport(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	today := utils.Today().Format("2006-01-02")

	var report models.MenuReport
	err = config.DB.Preload("Ingredients").Preload("Recommendations").
		Where("sppg_id = ? AND report_date = ?", userID, today).
		First(&report).Error
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":    "success",
			"submitted": false,
			"date":      today,
			"data":      nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"submitted": true,
		"date":      today,
		"data":      report,
	})
}

func SubmitMenuReport(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	type ingredientInput struct {
		Name string `json:"name"`
	}
	type recommendationInput struct {
		Allergen           string `json:"allergen"`
		OriginalIngredient string `json:"original_ingredient"`
		Substitute         string `json:"substitute"`
		Category           string `json:"category"`
		Reason             string `json:"reason"`
	}
	type submitBody struct {
		MenuName           string                `json:"menu_name"`
		Kalori             float64               `json:"kalori"`
		Protein            float64               `json:"protein"`
		Karbohidrat        float64               `json:"karbohidrat"`
		Lemak              float64               `json:"lemak"`
		TotalPorsi         int                   `json:"total_porsi"`
		KaloriPerPorsi     float64               `json:"kalori_per_porsi"`
		StatusGizi         string                `json:"status_gizi"`
		PhotoRef           string                `json:"photo_ref"`
		Ingredients        []ingredientInput     `json:"ingredients"`
		UsedRecommendation bool                  `json:"used_recommendation"`
		Recommendations    []recommendationInput `json:"recommendations"`
	}

	var body submitBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format body tidak valid."})
	}

	if body.TotalPorsi <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "total_porsi harus bernilai positif."})
	}

	today := utils.Today()
	todayStr := today.Format("2006-01-02")

	// Satu laporan per SPPG per hari. Dicek lebih dulu supaya pesannya jelas;
	// unique index di tabel yang menjadi penjaga sebenarnya kalau ada dua
	// request bersamaan.
	var existing models.MenuReport
	if config.DB.Where("sppg_id = ? AND report_date = ?", userID, todayStr).First(&existing).Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "Laporan menu hari ini sudah dikirim.",
			"data":    existing,
		})
	}

	// Status alergen dihitung ulang dari data alergi siswa, bukan diambil dari
	// kiriman klien — klien hanya menyetor nama bahannya.
	names := make([]string, 0, len(body.Ingredients))
	for _, ing := range body.Ingredients {
		if n := strings.TrimSpace(ing.Name); n != "" {
			names = append(names, n)
		}
	}
	profiles, schoolIDs := sppgSchoolIDs(userID)
	statuses, allergenNames, _ := classifyIngredients(names, schoolIDs)
	containsAllergen := len(allergenNames) > 0

	// "Pakai rekomendasi" hanya bermakna kalau memang ada alergen dan
	// penggantinya disertakan.
	usedRecommendation := body.UsedRecommendation && containsAllergen && len(body.Recommendations) > 0

	report := models.MenuReport{
		SppgID:             userID,
		MenuName:           strings.TrimSpace(body.MenuName),
		Kalori:             body.Kalori,
		Protein:            body.Protein,
		Karbohidrat:        body.Karbohidrat,
		Lemak:              body.Lemak,
		TotalPorsi:         body.TotalPorsi,
		KaloriPerPorsi:     body.KaloriPerPorsi,
		StatusGizi:         body.StatusGizi,
		ReportDate:         today,
		ContainsAllergen:   containsAllergen,
		UsedRecommendation: usedRecommendation,
		PhotoPath:          resolveScanPhoto(body.PhotoRef),
	}

	tx := config.DB.Begin()
	if err := tx.Create(&report).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan laporan menu."})
	}

	for _, s := range statuses {
		row := models.MenuReportIngredient{
			MenuReportID: report.ID,
			Name:         s.Name,
			IsAllergen:   s.IsAllergen,
		}
		if err := tx.Create(&row).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan bahan terdeteksi."})
		}
		report.Ingredients = append(report.Ingredients, row)
	}

	if usedRecommendation {
		for _, r := range body.Recommendations {
			row := models.MenuReportRecommendation{
				MenuReportID:       report.ID,
				Allergen:           r.Allergen,
				OriginalIngredient: r.OriginalIngredient,
				Substitute:         r.Substitute,
				Category:           r.Category,
				Reason:             r.Reason,
			}
			if err := tx.Create(&row).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan rekomendasi menu."})
			}
			report.Recommendations = append(report.Recommendations, row)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan laporan menu."})
	}

	notifyMenuReport(userID, report, allergenNames, profiles)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Laporan menu berhasil dikirim.",
		"data":    report,
	})
}

// notifyMenuReport memberi tahu SPPG bahwa laporannya terkirim, dan — kalau
// menu mengandung alergen yang belum ditangani — memperingatkan sekolah yang
// siswanya terdampak.
func notifyMenuReport(sppgID uint, report models.MenuReport, allergens []string, profiles []models.SchoolProfile) {
	var sppgUser models.User
	config.DB.Select("id, name").First(&sppgUser, sppgID)

	services.Notify(
		sppgID,
		models.NotifLaporanTerkirim,
		"Laporan Terkirim!",
		fmt.Sprintf("Laporan menu %s pada tanggal %s berhasil terkirim!",
			defaultMenuName(report.MenuName), utils.TanggalID(report.ReportDate)),
		"menu_report", report.ID,
	)

	if !report.ContainsAllergen || report.UsedRecommendation || len(profiles) == 0 {
		return
	}

	schoolIDs := make([]uint, 0, len(profiles))
	for _, p := range profiles {
		schoolIDs = append(schoolIDs, p.UserID)
	}

	services.NotifyMany(
		schoolIDs,
		models.NotifPeringatanAlergen,
		"Peringatan Alergi",
		fmt.Sprintf("Menu %s dari %s tanggal %s mengandung %s. Periksa siswa dengan alergi tersebut.",
			defaultMenuName(report.MenuName), sppgUser.Name,
			utils.TanggalID(report.ReportDate), strings.Join(allergens, ", ")),
		"menu_report", report.ID,
	)
}

func defaultMenuName(name string) string {
	if strings.TrimSpace(name) == "" {
		return "hari ini"
	}
	return name
}
