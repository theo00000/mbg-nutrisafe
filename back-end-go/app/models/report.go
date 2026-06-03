package models

import "time"

type FoodReport struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SchoolID    uint      `gorm:"not null" json:"school_id"`
	SppgID      uint      `gorm:"not null" json:"sppg_id"`
	IssueType   string    `gorm:"type:varchar(50);not null" json:"issue_type"` // food_shortage, spoiled_food, allergen, late_delivery
	Description string    `gorm:"type:text" json:"description"`
	PhotoURL    string    `json:"photo_url"`
	ReportDate  time.Time `gorm:"type:date;not null" json:"report_date"`
	Status      string    `gorm:"type:varchar(30);default:'belum_ditinjau'" json:"status"` // belum_ditinjau, sedang_ditinjau, selesai
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

