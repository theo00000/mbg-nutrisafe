package models

import "time"

type SppgFoodProblem struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	SppgID            uint      `gorm:"not null;index" json:"sppg_id"`
	SchoolName        string    `gorm:"type:varchar(150);not null" json:"school_name"`
	FoodName          string    `gorm:"type:varchar(150);not null" json:"food_name"`
	ProblemType       string    `gorm:"type:varchar(80);not null" json:"problem_type"`
	AffectedPortions  int       `gorm:"not null" json:"affected_portions"`
	Priority          string    `gorm:"type:varchar(20);not null;default:'Sedang'" json:"priority"`
	Description       string    `gorm:"type:text;not null" json:"description"`
	PhotoPath         string    `gorm:"type:varchar(255)" json:"photo_path"`
	Status            string    `gorm:"type:varchar(30);default:'baru'" json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
