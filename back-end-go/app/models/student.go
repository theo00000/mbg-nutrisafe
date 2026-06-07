package models

import "time"

type Student struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	SchoolID         uint      `gorm:"not null;index" json:"school_id"`
	Name             string    `gorm:"type:varchar(100);not null" json:"name"`
	NISN             string    `gorm:"type:varchar(20);unique;not null" json:"nisn"`
	Class            string    `gorm:"type:varchar(20)" json:"class"`
	Gender           string    `gorm:"type:varchar(20)" json:"gender"`
	Address          string    `gorm:"type:text" json:"address"`
	AccountGenerated bool      `gorm:"default:false" json:"account_generated"`
	UserID           *uint     `json:"-"`
	ParentID         *uint     `json:"-"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	Allergies []Allergy `gorm:"many2many:student_allergies;" json:"allergies"`
}

type Allergy struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `json:"description"`
}