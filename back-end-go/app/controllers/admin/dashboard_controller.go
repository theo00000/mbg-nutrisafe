package admin

import (
	"back-end/app/models"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
)

func GetAdminDashboard(c *fiber.Ctx) error {
	var schoolCount, studentCount, totalsppg, totalUmum int64
	var sppgActive, sppgPending, sppgInactive int64

	if err := config.DB.Model(&models.SchoolProfile{}).Count(&schoolCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data dashboard. Silakan coba lagi nanti.",
		})
	}

	config.DB.Model(&models.Student{}).Count(&studentCount)
	config.DB.Model(&models.SppgProfile{}).Count(&totalsppg)
	config.DB.Model(&models.UmumProfile{}).Count(&totalUmum)
	config.DB.Model(&models.SppgProfile{}).Where("status = ?", "active").Count(&sppgActive)
	config.DB.Model(&models.SppgRegistration{}).Where("status IN ?", []string{"pending", "belum_ditinjau", "verifikasi"}).Count(&sppgPending)
	config.DB.Model(&models.SppgProfile{}).Where("status = ?", "inactive").Count(&sppgInactive)

	var recentsppg []models.RecentSppgEntry
	config.DB.Table("sppg_profiles").
		Select("sppg_profiles.user_id, users.name, users.email, users.phone, sppg_profiles.institution_name, sppg_profiles.capacity, sppg_profiles.address, sppg_profiles.status, users.created_at").
		Joins("JOIN users ON users.id = sppg_profiles.user_id").
		Order("users.created_at DESC").
		Limit(3).
		Scan(&recentsppg)

	var recentUmum []models.RecentUmumEntry
	config.DB.Table("umum_profiles").
		Select("umum_profiles.user_id, users.name, users.email, users.phone, umum_profiles.lokasi, users.created_at").
		Joins("JOIN users ON users.id = umum_profiles.user_id").
		Order("users.created_at DESC").
		Limit(5).
		Scan(&recentUmum)

	var recentSchool []models.RecentSchoolEntry
	config.DB.Table("school_profiles").
		Select("school_profiles.user_id, users.name, users.email, users.phone, school_profiles.npsn, school_profiles.grade, school_profiles.teacher_count, school_profiles.student_count, school_profiles.address, school_profiles.status, users.created_at").
		Joins("JOIN users ON users.id = school_profiles.user_id").
		Order("users.created_at DESC").
		Limit(5).
		Scan(&recentSchool)

	var recentReports []models.FoodReport
	config.DB.Order("created_at desc").Limit(5).Find(&recentReports)

	if recentsppg == nil {
		recentsppg = []models.RecentSppgEntry{}
	}
	if recentUmum == nil {
		recentUmum = []models.RecentUmumEntry{}
	}
	if recentSchool == nil {
		recentSchool = []models.RecentSchoolEntry{}
	}
	if recentReports == nil {
		recentReports = []models.FoodReport{}
	}

	return c.Status(fiber.StatusOK).JSON(models.DashboardResponse{
		Status: "success",
		Data: models.DashboardData{
			RecentReports: recentReports,
			Sppg:          recentsppg,
			SppgStatus: models.DashboardSppgStatus{
				Active:   sppgActive,
				Pending:  sppgPending,
				Inactive: sppgInactive,
			},
			Umum:         recentUmum,
			School:       recentSchool,
			TotalSchool:  schoolCount,
			TotalSppg:    totalsppg,
			TotalStudent: studentCount,
			TotalUmum:    totalUmum,
		},
	})
}
