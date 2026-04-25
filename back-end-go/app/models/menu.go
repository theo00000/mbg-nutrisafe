package models

import "time"

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