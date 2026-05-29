package models

import "time"

type MenuReport struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	SppgID         uint      `gorm:"not null" json:"sppg_id"`
	MenuName       string    `gorm:"type:varchar(200)" json:"menu_name"`
	Kalori         float64   `json:"kalori"`
	Protein        float64   `json:"protein"`
	Karbohidrat    float64   `json:"karbohidrat"`
	Lemak          float64   `json:"lemak"`
	TotalPorsi     int       `json:"total_porsi"`
	KaloriPerPorsi float64   `json:"kalori_per_porsi"`
	StatusGizi     string    `gorm:"type:varchar(50)" json:"status_gizi"`
	ReportDate     time.Time `gorm:"type:date;not null" json:"report_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Menu struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	Name         string       `gorm:"type:varchar(100);not null" json:"name"`
	Description  string       `json:"description"`
	CalorieCount int          `json:"calorie_count"`
	CreatedBy    uint         `json:"created_by"` 
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	
	Ingredients []Ingredient `gorm:"many2many:menu_ingredients;" json:"ingredients"`
}

type Ingredient struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `json:"description"`
}