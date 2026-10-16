package models

import (
	"encoding/json"
	"time"

	"back-end/utils"
)

type FoodReport struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SchoolID    uint      `gorm:"not null" json:"school_id"`
	SppgID      uint      `gorm:"not null" json:"sppg_id"`
	ReporterID  uint      `gorm:"index" json:"reporter_id"` // user pengirim (sekolah / umum / siswa)
	IssueType   string    `gorm:"type:varchar(50);not null" json:"issue_type"` // food_shortage, spoiled_food, allergen, late_delivery
	Description string    `gorm:"type:text" json:"description"`
	PhotoURL    string    `json:"photo_url"`
	ReportDate  time.Time `gorm:"type:date;not null" json:"report_date"`
	Status      string    `gorm:"type:varchar(30);default:'belum_ditinjau'" json:"status"` // belum_ditinjau, sedang_ditinjau, selesai
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MarshalJSON menukar PhotoURL dari path disk ke URL publik. Ditaruh di model
// supaya setiap controller yang mengembalikan FoodReport apa adanya ikut benar,
// tanpa harus memetakan ulang satu per satu.
func (r FoodReport) MarshalJSON() ([]byte, error) {
	type alias FoodReport // alias memutus rekursi: alias tidak punya MarshalJSON
	out := alias(r)
	out.PhotoURL = utils.PublicURL(r.PhotoURL)
	return json.Marshal(out)
}
