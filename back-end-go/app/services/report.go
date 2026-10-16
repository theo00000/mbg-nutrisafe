package services

import (
	"fmt"
	"time"

	"back-end/app/models"
	"back-end/config"
	"back-end/utils"
)

// Label manusiawi untuk kode yang disimpan di basis data. Dipakai supaya
// layar detail laporan tidak perlu menerjemahkan sendiri.
var issueTypeLabels = map[string]string{
	"kekurangan_makanan": "Kekurangan Makanan",
	"makanan_basi":       "Makanan Basi",
	"masalah_alergen":    "Masalah Alergen",
	"terlambat":          "Keterlambatan Pengiriman",
	"lainnya":            "Lainnya",
}

var reportStatusLabels = map[string]string{
	"belum_ditinjau":  "Belum Ditinjau",
	"sedang_ditinjau": "Sedang Ditinjau",
	"selesai":         "Selesai",
}

func IssueTypeLabel(code string) string {
	if label, ok := issueTypeLabels[code]; ok {
		return label
	}
	return code
}

func ReportStatusLabel(code string) string {
	if label, ok := reportStatusLabels[code]; ok {
		return label
	}
	return code
}

// FoodReportDetail adalah bentuk laporan yang siap ditampilkan di layar
// "Detail Laporan": sudah membawa nama sekolah/SPPG dan label berbahasa
// Indonesia, bukan sekadar ID dan kode.
type FoodReportDetail struct {
	ID               uint      `json:"id"`
	SchoolID         uint      `json:"school_id"`
	SchoolName       string    `json:"school_name"`
	SppgID           uint      `json:"sppg_id"`
	SppgName         string    `json:"sppg_name"`
	IssueType        string    `json:"issue_type"`
	IssueTypeLabel   string    `json:"issue_type_label"`
	Description      string    `json:"description"`
	PhotoURL         string    `json:"photo_url"`
	ReportDate       string    `json:"report_date"`
	ReportDateLabel  string    `json:"report_date_label"`
	Status           string    `json:"status"`
	StatusLabel      string    `json:"status_label"`
	SubmittedLabel   string    `json:"submitted_label"`
	CreatedAt        time.Time `json:"created_at"`
}

// submittedLabel meniru teks kecil di pojok layar detail: laporan hari ini
// ditulis "Hari ini pada 12:00", yang lebih lama pakai tanggal.
func submittedLabel(t time.Time) string {
	now := time.Now()
	y1, m1, d1 := t.Date()
	y2, m2, d2 := now.Date()
	if y1 == y2 && m1 == m2 && d1 == d2 {
		return fmt.Sprintf("Hari ini pada %s", t.Format("15:04"))
	}
	return fmt.Sprintf("%s pada %s", utils.TanggalID(t), t.Format("15:04"))
}

// BuildFoodReportDetail melengkapi satu laporan dengan nama pihak terkait.
func BuildFoodReportDetail(r models.FoodReport) FoodReportDetail {
	var schoolName, sppgName string
	if r.SchoolID != 0 {
		var school models.User
		if config.DB.Select("id, name").First(&school, r.SchoolID).Error == nil {
			schoolName = school.Name
		}
	}
	if r.SppgID != 0 {
		var sppg models.User
		if config.DB.Select("id, name").First(&sppg, r.SppgID).Error == nil {
			sppgName = sppg.Name
		}
	}

	return FoodReportDetail{
		ID:              r.ID,
		SchoolID:        r.SchoolID,
		SchoolName:      schoolName,
		SppgID:          r.SppgID,
		SppgName:        sppgName,
		IssueType:       r.IssueType,
		IssueTypeLabel:  IssueTypeLabel(r.IssueType),
		Description:     r.Description,
		PhotoURL:        utils.PublicURL(r.PhotoURL),
		ReportDate:      r.ReportDate.Format("2006-01-02"),
		ReportDateLabel: utils.TanggalID(r.ReportDate),
		Status:          r.Status,
		StatusLabel:     ReportStatusLabel(r.Status),
		SubmittedLabel:  submittedLabel(r.CreatedAt),
		CreatedAt:       r.CreatedAt,
	}
}

// BuildFoodReportDetails memetakan banyak laporan sekaligus. Nama sekolah dan
// SPPG diambil dalam satu query, bukan dua query per baris.
func BuildFoodReportDetails(reports []models.FoodReport) []FoodReportDetail {
	out := make([]FoodReportDetail, 0, len(reports))
	if len(reports) == 0 {
		return out
	}

	idSet := make(map[uint]bool)
	for _, r := range reports {
		if r.SchoolID != 0 {
			idSet[r.SchoolID] = true
		}
		if r.SppgID != 0 {
			idSet[r.SppgID] = true
		}
	}
	ids := make([]uint, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	var users []models.User
	config.DB.Select("id, name").Where("id IN ?", ids).Find(&users)
	nameByID := make(map[uint]string, len(users))
	for _, u := range users {
		nameByID[u.ID] = u.Name
	}

	for _, r := range reports {
		out = append(out, FoodReportDetail{
			ID:              r.ID,
			SchoolID:        r.SchoolID,
			SchoolName:      nameByID[r.SchoolID],
			SppgID:          r.SppgID,
			SppgName:        nameByID[r.SppgID],
			IssueType:       r.IssueType,
			IssueTypeLabel:  IssueTypeLabel(r.IssueType),
			Description:     r.Description,
			PhotoURL:        utils.PublicURL(r.PhotoURL),
			ReportDate:      r.ReportDate.Format("2006-01-02"),
			ReportDateLabel: utils.TanggalID(r.ReportDate),
			Status:          r.Status,
			StatusLabel:     ReportStatusLabel(r.Status),
			SubmittedLabel:  submittedLabel(r.CreatedAt),
			CreatedAt:       r.CreatedAt,
		})
	}
	return out
}
