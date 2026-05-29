package cron

import (
	"fmt"
	"log"
	"time"

	"back-end/app/models"
	"back-end/config"
	"back-end/utils"

	"golang.org/x/crypto/bcrypt"
)

func StartStudentAccountCron() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		// jalankan sekali saat startup
		syncAllStudentCounts()
		processStudentAccounts()
		for range ticker.C {
			syncAllStudentCounts()
			processStudentAccounts()
		}
	}()
	log.Println("[cron] Student cron started (interval: 5 menit)")
}

// syncAllStudentCounts menyinkronkan student_count di semua SchoolProfile
// berdasarkan COUNT aktual dari tabel students, berjalan setiap 5 menit.
func syncAllStudentCounts() {
	var profiles []models.SchoolProfile
	if err := config.DB.Find(&profiles).Error; err != nil {
		log.Printf("[cron] syncAllStudentCounts: gagal query school_profiles: %v", err)
		return
	}

	updated := 0
	for _, profile := range profiles {
		var count int64
		config.DB.Model(&models.Student{}).Where("school_id = ?", profile.UserID).Count(&count)
		if int(count) != profile.StudentCount {
			config.DB.Model(&profile).Update("student_count", int(count))
			updated++
		}
	}

	if updated > 0 {
		log.Printf("[cron] syncAllStudentCounts: %d sekolah diperbarui", updated)
	}
}

func processStudentAccounts() {
	cutoff := time.Now().Add(-15 * time.Minute)

	var students []models.Student
	if err := config.DB.Where("account_generated = false AND updated_at <= ?", cutoff).Find(&students).Error; err != nil {
		log.Printf("[cron] gagal query siswa: %v", err)
		return
	}

	if len(students) == 0 {
		return
	}

	bySchool := map[uint][]models.Student{}
	for _, s := range students {
		bySchool[s.SchoolID] = append(bySchool[s.SchoolID], s)
	}

	for schoolID, schoolStudents := range bySchool {
		processSchoolStudents(schoolID, schoolStudents)
	}
}

func processSchoolStudents(schoolID uint, students []models.Student) {
	var schoolProfile models.SchoolProfile
	if err := config.DB.Where("user_id = ?", schoolID).First(&schoolProfile).Error; err != nil {
		log.Printf("[cron] profil sekolah %d tidak ditemukan: %v", schoolID, err)
		return
	}

	var schoolUser models.User
	if err := config.DB.First(&schoolUser, schoolID).Error; err != nil {
		log.Printf("[cron] user sekolah %d tidak ditemukan: %v", schoolID, err)
		return
	}

	var accounts []utils.StudentAccountInfo

	for _, student := range students {
		loginEmail := buildStudentEmail(student.Name, student.NISN)

		var existing models.User
		if config.DB.Where("email = ?", loginEmail).First(&existing).Error == nil {
			log.Printf("[cron] email %s sudah ada, skip siswa %s", loginEmail, student.Name)
			continue
		}

		tempPassword := utils.GeneratePasswordFromName(student.Name)
		hash, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("[cron] gagal hash password siswa %s: %v", student.Name, err)
			continue
		}

		user := models.User{
			Name:         student.Name,
			Email:        loginEmail,
			Phone:        "",
			PasswordHash: string(hash),
			RoleName:     "siswa",
		}
		if err := config.DB.Create(&user).Error; err != nil {
			log.Printf("[cron] gagal buat akun siswa %s: %v", student.Name, err)
			continue
		}

		userID := user.ID
		config.DB.Model(&student).Updates(map[string]interface{}{
			"account_generated": true,
			"user_id":           userID,
		})

		accounts = append(accounts, utils.StudentAccountInfo{
			Name:       student.Name,
			LoginEmail: loginEmail,
			Password:   tempPassword,
		})
	}

	if len(accounts) == 0 {
		return
	}

	if schoolProfile.ContactEmail != "" {
		if err := utils.SendStudentAccountList(schoolProfile.ContactEmail, schoolUser.Name, accounts); err != nil {
			log.Printf("[cron] gagal kirim email ke sekolah %s: %v", schoolProfile.ContactEmail, err)
		} else {
			log.Printf("[cron] berhasil kirim %d akun siswa ke %s", len(accounts), schoolProfile.ContactEmail)
		}
	}
}

func buildStudentEmail(name, nisn string) string {
	slug := utils.GenerateSlug(name)
	nisnSuffix := nisn
	if len(nisnSuffix) >= 4 {
		nisnSuffix = nisnSuffix[len(nisnSuffix)-4:]
	}
	return fmt.Sprintf("%s.%s@siswa.sch.id", slug, nisnSuffix)
}
