package models

import "time"

type AllergyAlternativeSuggestion struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	StudentID       uint      `json:"student_id"`
	AllergenID      uint      `json:"allergen_id"`       // FK merujuk ke alergi
	SuggestedMenuID uint      `json:"suggested_menu_id"` // FK merujuk ke menu
	SuggestionDate  time.Time `gorm:"type:date" json:"suggestion_date"`
	MLModelAccuracy float64   `json:"ml_model_accuracy"`
}