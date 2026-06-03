package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"back-end/app/models"
	"back-end/config"
	"back-end/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Name               string `json:"name"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
	Password           string `json:"password"`
	RoleName           string `json:"role_name"`
	RegistrationSecret string `json:"registration_secret"`
}

func Register(c *fiber.Ctx) error {
	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format request tidak valid. Pastikan mengirim JSON yang benar.",
		})
	}

	if input.Name == "" || input.Email == "" || input.Password == "" || input.RoleName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Semua kolom (name, email, password, role_name) wajib diisi!",
		})
	}

	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "Email ini sudah terdaftar. Silakan gunakan email lain atau langsung login.",
		})
	}

	if err := config.DB.Where("phone = ?", input.Phone).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "Nomor telepon ini sudah terdaftar. Silakan gunakan nomor lain.",
		})
	}

	if input.RoleName != "school" && input.RoleName != "umum" && input.RoleName != "admin" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Role tidak valid! Pilihan yang tersedia: school, umum.",
		})
	}

	if input.RoleName == "admin" {
		expected := os.Getenv("ADMIN_REGISTRATION_SECRET")
		if expected == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Registrasi admin saat ini dinonaktifkan.",
			})
		}
		headerSecret := c.Get("X-Registration-Secret")
		if input.RegistrationSecret != expected && headerSecret != expected {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Kode registrasi admin tidak valid.",
			})
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengenkripsi kata sandi.",
		})
	}

	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		Phone:        input.Phone,
		PasswordHash: string(hashedPassword),
		RoleName:     input.RoleName,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan data pengguna ke server. Silakan coba lagi nanti.",
		})
	}

	switch user.RoleName {
	case "school":
		config.DB.Create(&models.SchoolProfile{UserID: user.ID, Status: "active"})
	case "umum":
		config.DB.Create(&models.UmumProfile{UserID: user.ID})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Registrasi berhasil",
		"data": fiber.Map{
			"user_id": user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"phone":   user.Phone,
			"role":    user.RoleName,
		},
	})
}

func Registersppg(c *fiber.Ctx) error {
	registrantName := c.FormValue("registrant_name")
	institutionName := c.FormValue("institution_name")
	nikNpwp := c.FormValue("nik_npwp")
	email := c.FormValue("email")
	phone := c.FormValue("phone")
	sppgName := c.FormValue("sppg_name")
	sppgAddress := c.FormValue("sppg_address")
	productionCapacityStr := c.FormValue("production_capacity")

	if registrantName == "" || nikNpwp == "" || email == "" || phone == "" || sppgName == "" || sppgAddress == "" || productionCapacityStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Semua kolom wajib (kecuali institution_name) harus diisi.",
		})
	}

	productionCapacity, err := strconv.Atoi(productionCapacityStr)
	if err != nil || productionCapacity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Kapasitas produksi harus berupa angka positif.",
		})
	}

	if len(nikNpwp) != 16 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "NIK harus terdiri dari 16 digit.",
		})
	}

	var existing models.SppgRegistration
	if err := config.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "Email ini sudah pernah mendaftar sebagai mitra sppg.",
		})
	}

	hashedNikNpwp, err := bcrypt.GenerateFromPassword([]byte(nikNpwp), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengenkripsi data identitas.",
		})
	}

	uploadDir := "./uploads/sppg"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyiapkan direktori upload.",
		})
	}

	proposalFile, err := c.FormFile("proposal")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "File proposal wajib diupload.",
		})
	}
	if filepath.Ext(proposalFile.Filename) != ".pdf" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "File proposal harus berformat .pdf.",
		})
	}
	proposalFilename := fmt.Sprintf("proposal_%d_%s", time.Now().UnixNano(), filepath.Base(proposalFile.Filename))
	proposalPath := filepath.Join(uploadDir, proposalFilename)
	if err := c.SaveFile(proposalFile, proposalPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan file proposal.",
		})
	}

	kitchenPhotoFile, err := c.FormFile("kitchen_photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "File foto dapur wajib diupload.",
		})
	}
	if kitchenPhotoFile.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Ukuran foto dapur tidak boleh melebihi 5 MB.",
		})
	}
	kitchenPhotoFilename := fmt.Sprintf("kitchen_photo_%d_%s", time.Now().UnixNano(), filepath.Base(kitchenPhotoFile.Filename))
	kitchenPhotoPath := filepath.Join(uploadDir, kitchenPhotoFilename)
	if err := c.SaveFile(kitchenPhotoFile, kitchenPhotoPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan file foto dapur.",
		})
	}

	registration := models.SppgRegistration{
		RegistrantName:     registrantName,
		InstitutionName:    institutionName,
		NikNpwp:            string(hashedNikNpwp),
		Email:              email,
		Phone:              phone,
		SppgName:           sppgName,
		SppgAddress:        sppgAddress,
		ProductionCapacity: productionCapacity,
		ProposalPath:       proposalPath,
		KitchenPhotoPath:   kitchenPhotoPath,
		Status:             "pending",
	}

	if err := config.DB.Create(&registration).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan data pendaftaran. Silakan coba lagi nanti.",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Formulir pendaftaran berhasil dikirim. Silahkan tunggu verifikasi selanjutnya melalui e-mail.",
		"data": fiber.Map{
			"id":    registration.ID,
			"email": registration.Email,
		},
	})
}

type RegisterSchoolInput struct {
	Name       string `json:"name"`        // nama guru/penanggung jawab
	Email      string `json:"email"`       // email pribadi guru (untuk menerima credentials)
	Phone      string `json:"phone"`
	SchoolName string `json:"school_name"` // nama sekolah
	NPSN       string `json:"npsn"`
	Address    string `json:"address"`
	Grade      string `json:"grade"`       // SD/MI, SMP/MTs/MTsN, SMA/SMK/MA/MAN
}

func RegisterSchool(c *fiber.Ctx) error {
	var input RegisterSchoolInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format request tidak valid. Pastikan mengirim JSON yang benar.",
		})
	}

	if input.Name == "" || input.Email == "" || input.Phone == "" || input.SchoolName == "" || input.NPSN == "" || input.Address == "" || input.Grade == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Semua kolom (name, email, phone, school_name, npsn, address, grade) wajib diisi.",
		})
	}

	validGrades := map[string]bool{"SD/MI": true, "SMP/MTs/MTsN": true, "SMA/SMK/MA/MAN": true}
	if !validGrades[input.Grade] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Tingkat sekolah tidak valid.",
		})
	}

	// cek email pribadi guru belum dipakai di sekolah lain
	var existingProfile models.SchoolProfile
	if config.DB.Where("contact_email = ?", input.Email).First(&existingProfile).Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "Email ini sudah terdaftar untuk sekolah lain.",
		})
	}

	// cek NPSN belum dipakai
	if config.DB.Where("npsn = ?", input.NPSN).First(&existingProfile).Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "NPSN ini sudah terdaftar.",
		})
	}

	// generate login email dari nama sekolah
	slug := utils.GenerateSlug(input.SchoolName)
	loginEmail := slug + "@sch.id"
	counter := 2
	for {
		var existingUser models.User
		if config.DB.Where("email = ?", loginEmail).First(&existingUser).Error != nil {
			break
		}
		loginEmail = fmt.Sprintf("%s.%d@sch.id", slug, counter)
		counter++
	}

	// generate random password
	tempPassword := utils.GeneratePasswordFromName(input.SchoolName)
	hash, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengenkripsi kata sandi.",
		})
	}

	user := models.User{
		Name:         input.SchoolName,
		Email:        loginEmail,
		Phone:        input.Phone,
		PasswordHash: string(hash),
		RoleName:     "school",
	}
	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal membuat akun sekolah.",
		})
	}

	profile := models.SchoolProfile{
		UserID:       user.ID,
		NPSN:         input.NPSN,
		Grade:        input.Grade,
		Address:      input.Address,
		ContactEmail: input.Email,
		Status:       "active",
	}
	if err := config.DB.Create(&profile).Error; err != nil {
		config.DB.Delete(&user)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan profil sekolah.",
		})
	}

	// kirim email credentials ke guru secara asinkron
	go func() {
		if err := utils.SendSchoolCredentials(input.Email, input.SchoolName, loginEmail, tempPassword); err != nil {
			fmt.Printf("[mailer] gagal kirim email ke %s: %v\n", input.Email, err)
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Pendaftaran sekolah berhasil. Informasi login telah dikirim ke email " + input.Email + ".",
		"data": fiber.Map{
			"school_name": input.SchoolName,
			"login_email": loginEmail,
			"contact_email": input.Email,
		},
	})
}

