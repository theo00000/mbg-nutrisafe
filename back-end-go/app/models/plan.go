package models

import "time"

type DailyMenuPlan struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `json:"user_id"`
	MenuID      uint      `json:"menu_id"`
	PlannedDate time.Time `gorm:"type:date" json:"planned_date"`
	CreatedAt   time.Time `json:"created_at"`
}