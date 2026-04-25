package models

import "time"

type FoodReport struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `json:"user_id"`
	MenuID     uint      `json:"menu_id"`
	ReportDate time.Time `gorm:"type:date" json:"report_date"`
	PhotoURL   string    `json:"photo_url"`
	Comment    string    `json:"comment"`
	Status     string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, approved, rejected
}