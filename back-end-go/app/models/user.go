package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Email        string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Phone        string    `gorm:"type:varchar(20);not null" json:"phone"`
	PasswordHash string    `gorm:"not null" json:"-"`
	RoleName     string    `gorm:"type:varchar(50);not null" json:"role_name"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type Role struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(50);unique;not null" json:"name"`
	Description string `json:"description"`
}

type SchoolProfile struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	UserID       uint   `gorm:"not null;uniqueIndex" json:"user_id"`
	NPSN         string `gorm:"type:varchar(20)" json:"npsn"`
	Grade        string `gorm:"type:varchar(20)" json:"grade"` // SD/MI, SMP/MTs/MTsN, SMA/SMK/MA/MAN
	TeacherCount int    `json:"teacher_count"`
	StudentCount int    `json:"student_count"`
	Address      string `gorm:"type:text" json:"address"`
	ContactEmail string `gorm:"type:varchar(100)" json:"contact_email"`
	SppgUserID   *uint  `json:"sppg_user_id"`
	Status       string `gorm:"type:varchar(20);default:'active'" json:"status"`
}

type SppgProfile struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	UserID          uint   `gorm:"not null;uniqueIndex" json:"user_id"`
	InstitutionName string `gorm:"type:varchar(100)" json:"institution_name"`
	Capacity        int    `json:"capacity"`
	Address         string `gorm:"type:text" json:"address"`
	ContactEmail    string `gorm:"type:varchar(100)" json:"contact_email"`
	Status          string `gorm:"type:varchar(20);default:'active'" json:"status"`
}

type SppgRegistration struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	RegistrantName     string    `gorm:"type:varchar(100);not null" json:"registrant_name"`
	InstitutionName    string    `gorm:"type:varchar(100)" json:"institution_name"`
	NikNpwp            string    `gorm:"type:varchar(100);not null" json:"-"`
	Email              string    `gorm:"type:varchar(100);not null" json:"email"`
	Phone              string    `gorm:"type:varchar(20);not null" json:"phone"`
	SppgName           string    `gorm:"type:varchar(100);not null" json:"sppg_name"`
	SppgAddress        string    `gorm:"type:text;not null" json:"sppg_address"`
	ProductionCapacity int       `gorm:"not null" json:"production_capacity"`
	ProposalPath       string    `gorm:"type:varchar(255);not null" json:"-"`
	KitchenPhotoPath   string    `gorm:"type:varchar(255);not null" json:"-"`
	Status             string    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type UmumProfile struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserID uint   `gorm:"not null;uniqueIndex" json:"user_id"`
	Lokasi string `gorm:"type:varchar(255)" json:"lokasi"`
}
