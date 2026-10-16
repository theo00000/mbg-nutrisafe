package models

import (
	"encoding/json"
	"time"

	"back-end/utils"
)

// Status alergen sebuah laporan menu. Nilai inilah yang menentukan warna &
// teks banner di layar Scan & Pelaporan Menu Makanan.
const (
	AllergenStatusBebas           = "bebas_alergen"           // tak ada bahan yang cocok data alergi siswa
	AllergenStatusMengandung      = "mengandung_alergen"      // ada alergen, belum ditangani
	AllergenStatusAmanRekomendasi = "aman_dengan_rekomendasi" // ada alergen, tapi SPPG memakai menu pengganti
)

type MenuReport struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	SppgID         uint    `gorm:"not null;uniqueIndex:idx_menu_report_harian" json:"sppg_id"`
	MenuName       string  `gorm:"type:varchar(200)" json:"menu_name"`
	Kalori         float64 `json:"kalori"`
	Protein        float64 `json:"protein"`
	Karbohidrat    float64 `json:"karbohidrat"`
	Lemak          float64 `json:"lemak"`
	TotalPorsi     int     `json:"total_porsi"`
	KaloriPerPorsi float64 `json:"kalori_per_porsi"`
	StatusGizi     string  `gorm:"type:varchar(50)" json:"status_gizi"`

	// Satu SPPG hanya boleh mengirim satu laporan menu per hari — layar
	// menampilkan "Laporan Hari Ini Telah di Kirim!" setelah terkirim.
	ReportDate time.Time `gorm:"type:date;not null;uniqueIndex:idx_menu_report_harian" json:"report_date"`

	ContainsAllergen   bool   `gorm:"default:false" json:"contains_allergen"`
	UsedRecommendation bool   `gorm:"default:false" json:"used_recommendation"`
	PhotoPath          string `gorm:"type:varchar(255)" json:"photo_path"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Ingredients     []MenuReportIngredient     `gorm:"foreignKey:MenuReportID" json:"ingredients"`
	Recommendations []MenuReportRecommendation `gorm:"foreignKey:MenuReportID" json:"recommendations"`
}

// AllergenStatus menurunkan status tampilan dari dua flag di atas, supaya
// klien tidak perlu menyusun logikanya sendiri.
func (m MenuReport) AllergenStatus() string {
	switch {
	case !m.ContainsAllergen:
		return AllergenStatusBebas
	case m.UsedRecommendation:
		return AllergenStatusAmanRekomendasi
	default:
		return AllergenStatusMengandung
	}
}

// MarshalJSON menambahkan status turunan dan mengubah foto jadi URL publik.
func (m MenuReport) MarshalJSON() ([]byte, error) {
	type alias MenuReport
	out := struct {
		alias
		AllergenStatus string `json:"allergen_status"`
	}{alias: alias(m), AllergenStatus: m.AllergenStatus()}
	out.PhotoPath = utils.PublicURL(m.PhotoPath)
	return json.Marshal(out)
}

// MenuReportIngredient menyimpan satu bahan hasil deteksi beserta apakah bahan
// itu cocok dengan data alergi siswa di sekolah yang dilayani SPPG.
// Tanpa tabel ini, hasil scan hilang begitu respons ditutup.
type MenuReportIngredient struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	MenuReportID uint   `gorm:"not null;index" json:"menu_report_id"`
	Name         string `gorm:"type:varchar(100);not null" json:"name"`
	IsAllergen   bool   `gorm:"default:false" json:"is_allergen"`
}

// MenuReportRecommendation adalah usulan bahan pengganti dari ML service yang
// benar-benar dipakai SPPG untuk laporan tersebut.
type MenuReportRecommendation struct {
	ID                 uint   `gorm:"primaryKey" json:"id"`
	MenuReportID       uint   `gorm:"not null;index" json:"menu_report_id"`
	Allergen           string `gorm:"type:varchar(100);not null" json:"allergen"`
	OriginalIngredient string `gorm:"type:varchar(100)" json:"original_ingredient"`
	Substitute         string `gorm:"type:varchar(150);not null" json:"substitute"`
	Category           string `gorm:"type:varchar(80)" json:"category"`
	Reason             string `gorm:"type:text" json:"reason"`
}

type Menu struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Description  string    `json:"description"`
	CalorieCount int       `json:"calorie_count"`
	CreatedBy    uint      `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Ingredients []Ingredient `gorm:"many2many:menu_ingredients;" json:"ingredients"`
}

type Ingredient struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `json:"description"`
}
