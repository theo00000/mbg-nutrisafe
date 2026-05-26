package models

import "time"

type StudentAllergy struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	SchoolID       uint      `gorm:"not null" json:"school_id"` 
	StudentName    string    `gorm:"type:varchar(150);not null" json:"student_name"`
	ClassName      string    `gorm:"type:varchar(50);not null" json:"class_name"`
	AllergyType    string    `gorm:"type:varchar(100);not null" json:"allergy_type"`
	Description    string    `gorm:"type:text" json:"description"`
	Severity       string    `gorm:"type:varchar(50);not null" json:"severity"`
	ActionRequired string    `gorm:"type:text" json:"action_required"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}