package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Email        string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Phone        string    `gorm:"type:varchar(20);not null" json:"phone"`
	PasswordHash string    `gorm:"not null" json:"-"`
	RoleName     string    `gorm:"type:varchar(50);not null" json:"role_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Role struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(50);unique;not null" json:"name"`
	Description string `json:"description"`
}

type SchoolProfile struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	UserID       uint   `gorm:"not null;uniqueIndex" json:"user_id"`
	Level        string `gorm:"type:varchar(10)" json:"level"` // SD, SMP, SMA
	TeacherCount int    `json:"teacher_count"`
	StudentCount int    `json:"student_count"`
	Address      string `gorm:"type:text" json:"address"`
	SpggUserID   *uint  `json:"spgg_user_id"`
	Status       string `gorm:"type:varchar(20);default:'active'" json:"status"`
}

type SpggProfile struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   uint   `gorm:"not null;uniqueIndex" json:"user_id"`
	Capacity int    `json:"capacity"`
	Address  string `gorm:"type:text" json:"address"`
	Status   string `gorm:"type:varchar(20);default:'active'" json:"status"`
}

type SpggRegistration struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	RegistrantName     string    `gorm:"type:varchar(100);not null" json:"registrant_name"`
	InstitutionName    string    `gorm:"type:varchar(100)" json:"institution_name"`
	NikNpwp            string    `gorm:"type:varchar(30);not null" json:"nik_npwp"`
	Email              string    `gorm:"type:varchar(100);not null" json:"email"`
	Phone              string    `gorm:"type:varchar(20);not null" json:"phone"`
	SpggName           string    `gorm:"type:varchar(100);not null" json:"spgg_name"`
	SpggAddress        string    `gorm:"type:text;not null" json:"spgg_address"`
	ProductionCapacity int       `gorm:"not null" json:"production_capacity"`
	ProposalPath       string    `gorm:"type:varchar(255);not null" json:"proposal_path"`
	KitchenPhotoPath   string    `gorm:"type:varchar(255);not null" json:"kitchen_photo_path"`
	Status             string    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}