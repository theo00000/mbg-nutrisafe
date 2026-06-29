package school

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

var onlyDigits = regexp.MustCompile(`^\d+$`)

type StudentInput struct {
	Name    string `json:"name"`
	NISN    string `json:"nisn"`
	Class   string `json:"class"`
	Gender  string `json:"gender"`
	Address string `json:"address"`
}

type BulkStudentInput struct {
	Students []StudentInput `json:"students"`
}

type ImportResult struct {
	Row     int    `json:"row"`
	NISN    string `json:"nisn"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func ImportStudents(c *fiber.Ctx) error {
	schoolID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var input BulkStudentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format request tidak valid."})
	}

	if len(input.Students) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Tidak ada data siswa untuk diimport."})
	}

	if len(input.Students) > 1000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Maksimal 1000 data siswa per import."})
	}

	nisnList := make([]string, 0, len(input.Students))
	for _, s := range input.Students {
		if s.NISN != "" {
			nisnList = append(nisnList, s.NISN)
		}
	}
	existingNISNs := map[string]bool{}
	if len(nisnList) > 0 {
		var found []string
		config.DB.Model(&models.Student{}).Where("nisn IN ?", nisnList).Pluck("nisn", &found)
		for _, n := range found {
			existingNISNs[n] = true
		}
	}

	results := make([]ImportResult, 0, len(input.Students))
	toInsert := make([]models.Student, 0, len(input.Students))
	seenInBatch := map[string]bool{}
	successCount, skipCount, errorCount := 0, 0, 0

	for i, s := range input.Students {
		res := ImportResult{Row: i + 2, NISN: s.NISN, Name: s.Name}

		if s.Name == "" || s.NISN == "" || s.Class == "" || s.Gender == "" {
			res.Status = "error"
			res.Message = "Kolom Nama, NISN, Kelas, dan Jenis Kelamin wajib diisi."
			results = append(results, res)
			errorCount++
			continue
		}
		if len(s.NISN) != 10 || !onlyDigits.MatchString(s.NISN) {
			res.Status = "error"
			res.Message = "NISN harus terdiri dari 10 digit angka."
			results = append(results, res)
			errorCount++
			continue
		}
		if s.Gender != "Laki-laki" && s.Gender != "Perempuan" {
			res.Status = "error"
			res.Message = "Jenis kelamin harus 'Laki-laki' atau 'Perempuan'."
			results = append(results, res)
			errorCount++
			continue
		}
		if existingNISNs[s.NISN] || seenInBatch[s.NISN] {
			res.Status = "skipped"
			res.Message = "NISN sudah terdaftar."
			results = append(results, res)
			skipCount++
			continue
		}

		toInsert = append(toInsert, models.Student{
			SchoolID: schoolID,
			Name:     s.Name,
			NISN:     s.NISN,
			Class:    s.Class,
			Gender:   s.Gender,
			Address:  s.Address,
		})
		seenInBatch[s.NISN] = true

		res.Status = "success"
		results = append(results, res)
		successCount++
	}

	classesCreated := 0
	if len(toInsert) > 0 {
		classSet := map[string]bool{}
		classList := make([]string, 0)
		for _, s := range toInsert {
			if !classSet[s.Class] {
				classSet[s.Class] = true
				classList = append(classList, s.Class)
			}
		}
		var existingClasses []string
		config.DB.Model(&models.SchoolClass{}).
			Where("school_id = ? AND name IN ?", schoolID, classList).
			Pluck("name", &existingClasses)
		existingClassSet := map[string]bool{}
		for _, n := range existingClasses {
			existingClassSet[n] = true
		}
		newClasses := make([]models.SchoolClass, 0)
		for _, name := range classList {
			if !existingClassSet[name] {
				newClasses = append(newClasses, models.SchoolClass{SchoolID: schoolID, Name: name})
			}
		}
		if len(newClasses) > 0 {
			if err := config.DB.Create(&newClasses).Error; err == nil {
				classesCreated = len(newClasses)
			}
		}

		if err := config.DB.CreateInBatches(&toInsert, 100).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan data siswa ke database."})
		}
		syncStudentCount(schoolID)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Proses import selesai.",
		"summary": fiber.Map{
			"total":           len(input.Students),
			"success":         successCount,
			"skipped":         skipCount,
			"error":           errorCount,
			"classes_created": classesCreated,
		},
		"results": results,
	})
}

func AddStudent(c *fiber.Ctx) error {
	schoolID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var input StudentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format request tidak valid."})
	}

	if input.Name == "" || input.NISN == "" || input.Class == "" || input.Gender == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Kolom name, nisn, class, dan gender wajib diisi.",
		})
	}

	if len(input.NISN) != 10 || !onlyDigits.MatchString(input.NISN) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "NISN harus terdiri dari 10 digit angka."})
	}

	if input.Gender != "Laki-laki" && input.Gender != "Perempuan" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Jenis kelamin harus 'Laki-laki' atau 'Perempuan'."})
	}

	var existing models.Student
	if config.DB.Where("nisn = ?", input.NISN).First(&existing).Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "error", "message": "NISN ini sudah terdaftar."})
	}

	student := models.Student{
		SchoolID: schoolID,
		Name:     input.Name,
		NISN:     input.NISN,
		Class:    input.Class,
		Gender:   input.Gender,
		Address:  input.Address,
	}

	if err := config.DB.Create(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan data siswa."})
	}

	syncStudentCount(schoolID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Data siswa berhasil disimpan.",
		"data":    student,
	})
}

func syncStudentCount(schoolID uint) {
	var count int64
	config.DB.Model(&models.Student{}).Where("school_id = ?", schoolID).Count(&count)
	config.DB.Model(&models.SchoolProfile{}).Where("user_id = ?", schoolID).Update("student_count", count)
}

func UpdateStudent(c *fiber.Ctx) error {
	schoolID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	id := c.Params("id")
	var student models.Student
	if err := config.DB.Where("id = ? AND school_id = ?", id, schoolID).First(&student).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Data siswa tidak ditemukan."})
	}

	if student.AccountGenerated {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Data siswa tidak dapat diubah karena akun sudah dibuat.",
		})
	}

	var input StudentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format request tidak valid."})
	}

	if input.Gender != "" && input.Gender != "Laki-laki" && input.Gender != "Perempuan" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Jenis kelamin harus 'Laki-laki' atau 'Perempuan'."})
	}

	if input.Name != "" {
		student.Name = input.Name
	}
	if input.Class != "" {
		student.Class = input.Class
	}
	if input.Gender != "" {
		student.Gender = input.Gender
	}
	if input.Address != "" {
		student.Address = input.Address
	}

	if err := config.DB.Save(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal memperbarui data siswa."})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Data siswa berhasil diperbarui.",
		"data":    student,
	})
}

func GetStudents(c *fiber.Ctx) error {
	schoolID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var students []models.Student
	config.DB.Where("school_id = ?", schoolID).Find(&students)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   students,
	})
}

func DeleteStudent(c *fiber.Ctx) error {
	schoolID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	id := c.Params("id")
	var student models.Student
	if err := config.DB.Where("id = ? AND school_id = ?", id, schoolID).First(&student).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Data siswa tidak ditemukan."})
	}

	if student.AccountGenerated {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Data siswa tidak dapat dihapus karena akun sudah dibuat.",
		})
	}

	config.DB.Delete(&student)
	syncStudentCount(schoolID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Data siswa berhasil dihapus."})
}
