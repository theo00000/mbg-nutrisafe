package models

import "time"

type DeliveryReport struct {
	ID                  uint               `gorm:"primaryKey" json:"id"`
	SppgID              uint               `gorm:"not null" json:"sppg_id"`
	SchoolID            uint               `gorm:"not null" json:"school_id"`
	DeliveryDate        time.Time          `gorm:"type:date;not null" json:"delivery_date"`
	TotalPortions       int                `gorm:"not null" json:"total_portions"`
	DistributedPortions int                `json:"distributed_portions"`
	Status              string             `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, diterima
	Photo1Path          string             `gorm:"type:varchar(255)" json:"photo1_path"`
	Photo2Path          string             `gorm:"type:varchar(255)" json:"photo2_path"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
	Items               []DeliveryMenuItem `gorm:"foreignKey:DeliveryReportID" json:"items,omitempty"`
}

type DeliveryMenuItem struct {
	ID                  uint   `gorm:"primaryKey" json:"id"`
	DeliveryReportID    uint   `gorm:"not null" json:"delivery_report_id"`
	MenuName            string `gorm:"type:varchar(100);not null" json:"menu_name"`
	Category            string `gorm:"type:varchar(50)" json:"category"` // Karbohidrat, Protein Hewani, Sayuran, Buah, dll
	Portions            int    `gorm:"not null" json:"portions"`
	IsAllergySubstitute bool   `gorm:"default:false" json:"is_allergy_substitute"`
	AllergyNote         string `gorm:"type:varchar(100)" json:"allergy_note"` // contoh: "tanpa telur"
}
