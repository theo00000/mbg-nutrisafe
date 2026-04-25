package models

import "time"

type Student struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Class     string    `gorm:"type:varchar(20)" json:"class"`
	NISN      string    `gorm:"type:varchar(20);unique" json:"nisn"`
	ParentID  uint      `json:"parent_id"` // FK merujuk ke id di tabel users
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	Allergies []Allergy `gorm:"many2many:student_allergies;" json:"allergies"`
}

type Allergy struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `json:"description"`
}