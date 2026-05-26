package controllers

import (
	"fmt"
	"os"
	"strings"

	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type StatusUpdateInput struct {
	Status string `json:"status"`
}

func requireAdmin(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fmt.Errorf("token tidak ditemukan")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fmt.Errorf("format token tidak valid")
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return fmt.Errorf("konfigurasi keamanan server tidak lengkap")
	}
	token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode enkripsi tidak valid")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("token tidak valid atau sudah kadaluarsa")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("gagal membaca token")
	}
	role, _ := claims["role"].(string)
	if role != "admin" {
		return fmt.Errorf("hanya admin yang dapat mengakses fitur ini")
	}
	return nil
}

// ============================================================
// DASHBOARD
// ============================================================

func GetAdminDashboard(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var schoolCount, studentCount, totalSpgg, totalUmum int64
	var spggActive, spggPending, spggInactive int64

	if err := config.DB.Model(&models.User{}).Where("role_name = ?", "school").Count(&schoolCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data dashboard. Silakan coba lagi nanti.",
		})
	}

	config.DB.Model(&models.Student{}).Count(&studentCount)
	config.DB.Model(&models.User{}).Where("role_name = ?", "spgg").Count(&totalSpgg)
	config.DB.Model(&models.User{}).Where("role_name = ?", "umum").Count(&totalUmum)
	config.DB.Model(&models.SpggProfile{}).Where("status = ?", "active").Count(&spggActive)
	config.DB.Model(&models.SpggRegistration{}).Where("status = ?", "verifikasi").Count(&spggPending)
	config.DB.Model(&models.SpggProfile{}).Where("status = ?", "inactive").Count(&spggInactive)

	var recentRegistrations []models.SpggRegistration
	config.DB.Order("created_at desc").Limit(5).Find(&recentRegistrations)

	var recentReports []models.FoodReport
	config.DB.Order("created_at desc").Limit(5).Find(&recentReports)

	if recentRegistrations == nil {
		recentRegistrations = []models.SpggRegistration{}
	}
	if recentReports == nil {
		recentReports = []models.FoodReport{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"total_school":  schoolCount,
			"total_spgg":    totalSpgg,
			"total_student": studentCount,
			"total_umum":    totalUmum,
			"spgg_status": fiber.Map{
				"active":   spggActive,
				"pending":  spggPending,
				"inactive": spggInactive,
			},
			"recent_registrations": recentRegistrations,
			"recent_reports":       recentReports,
		},
	})
}

// ============================================================
// PENDAFTARAN MITRA
// ============================================================

var validRegistrationStatuses = map[string]bool{
	"belum_ditinjau": true,
	"verifikasi":     true,
	"disetujui":      true,
}

func GetRegistrations(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	status := c.Query("status")
	if status != "" && !validRegistrationStatuses[status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Filter status tidak valid. Pilihan: belum_ditinjau, verifikasi, disetujui.",
		})
	}

	var registrations []models.SpggRegistration
	query := config.DB.Order("created_at desc")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&registrations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data pendaftaran.",
		})
	}

	if registrations == nil {
		registrations = []models.SpggRegistration{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   registrations,
	})
}

func GetRegistrationDetail(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	id := c.Params("id")
	var registration models.SpggRegistration
	if err := config.DB.First(&registration, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Data pendaftaran tidak ditemukan.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   registration,
	})
}

func UpdateRegistrationStatus(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	id := c.Params("id")
	var input StatusUpdateInput
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

	if !validRegistrationStatuses[input.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Status tidak valid. Pilihan: belum_ditinjau, verifikasi, disetujui.",
		})
	}

	var registration models.SpggRegistration
	if err := config.DB.First(&registration, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Data pendaftaran tidak ditemukan.",
		})
	}

	if err := config.DB.Model(&registration).Update("status", input.Status).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal memperbarui status pendaftaran.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Status pendaftaran berhasil diperbarui.",
		"data":    registration,
	})
}

// ============================================================
// KELOLA SEKOLAH
// ============================================================

func GetSchools(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var schools []models.User
	if err := config.DB.Where("role_name = ?", "school").Find(&schools).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data sekolah.",
		})
	}

	type SchoolItem struct {
		ID           uint   `json:"id"`
		Name         string `json:"name"`
		Email        string `json:"email"`
		Level        string `json:"level"`
		Address      string `json:"address"`
		SpggName     string `json:"spgg_name"`
		TeacherCount int    `json:"teacher_count"`
		StudentCount int    `json:"student_count"`
		Status       string `json:"status"`
	}

	result := make([]SchoolItem, 0)
	for _, school := range schools {
		var profile models.SchoolProfile
		config.DB.Where("user_id = ?", school.ID).First(&profile)

		spggName := ""
		if profile.SpggUserID != nil {
			var spgg models.User
			if config.DB.First(&spgg, *profile.SpggUserID).Error == nil {
				spggName = spgg.Name
			}
		}

		result = append(result, SchoolItem{
			ID:           school.ID,
			Name:         school.Name,
			Email:        school.Email,
			Level:        profile.Level,
			Address:      profile.Address,
			SpggName:     spggName,
			TeacherCount: profile.TeacherCount,
			StudentCount: profile.StudentCount,
			Status:       profile.Status,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func GetSchoolDetail(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	id := c.Params("id")
	var school models.User
	if err := config.DB.Where("id = ? AND role_name = ?", id, "school").First(&school).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Sekolah tidak ditemukan.",
		})
	}

	var profile models.SchoolProfile
	config.DB.Where("user_id = ?", school.ID).First(&profile)

	spggName := ""
	if profile.SpggUserID != nil {
		var spgg models.User
		if config.DB.First(&spgg, *profile.SpggUserID).Error == nil {
			spggName = spgg.Name
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":            school.ID,
			"name":          school.Name,
			"email":         school.Email,
			"phone":         school.Phone,
			"level":         profile.Level,
			"address":       profile.Address,
			"teacher_count": profile.TeacherCount,
			"student_count": profile.StudentCount,
			"spgg_name":     spggName,
			"status":        profile.Status,
		},
	})
}

// ============================================================
// MANAJEMEN AKUN
// ============================================================

var validRoles = map[string]bool{
	"admin":  true,
	"school": true,
	"spgg":   true,
	"umum":   true,
}

func GetUsers(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	role := c.Query("role")
	if role != "" && !validRoles[role] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Filter role tidak valid. Pilihan: admin, school, spgg, umum.",
		})
	}

	var users []models.User
	query := config.DB
	if role != "" {
		query = query.Where("role_name = ?", role)
	}
	if err := query.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data pengguna.",
		})
	}

	if users == nil {
		users = []models.User{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   users,
	})
}

type AdminUpdateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func AdminUpdateUser(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	id := c.Params("id")
	var input AdminUpdateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format request tidak valid.",
		})
	}

	if input.Name == "" && input.Email == "" && input.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Tidak ada data yang diubah. Isi minimal satu field (name, email, atau phone).",
		})
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Pengguna tidak ditemukan.",
		})
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" && input.Email != user.Email {
		var check models.User
		if config.DB.Where("email = ?", input.Email).First(&check).Error == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status":  "error",
				"message": "Email ini sudah digunakan oleh pengguna lain.",
			})
		}
		user.Email = input.Email
	}
	if input.Phone != "" && input.Phone != user.Phone {
		var check models.User
		if config.DB.Where("phone = ?", input.Phone).First(&check).Error == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status":  "error",
				"message": "Nomor telepon ini sudah digunakan oleh pengguna lain.",
			})
		}
		user.Phone = input.Phone
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan perubahan data pengguna.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Data pengguna berhasil diperbarui.",
		"data":    user,
	})
}

// ============================================================
// LAPORAN MASUK
// ============================================================

var validReportStatuses = map[string]bool{
	"belum_ditinjau":  true,
	"sedang_ditinjau": true,
	"selesai":         true,
}

func GetReports(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

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
	if err := requireAdmin(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

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
	if err := requireAdmin(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	id := c.Params("id")
	var input StatusUpdateInput
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

// ============================================================
// DISTRIBUSI MAKANAN
// ============================================================

var validDistributionStatuses = map[string]bool{
	"dalam_pengiriman": true,
	"diterima":         true,
}

func GetDistributions(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

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
			"items": distributions,
		},
	})
}

func UpdateDistributionStatus(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	id := c.Params("id")
	var input StatusUpdateInput
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
