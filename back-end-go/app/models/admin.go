package models

import "time"

type StatusUpdateInput struct {
	Status string `json:"status"`
}

type DashboardSppgStatus struct {
	Active   int64 `json:"active"`
	Pending  int64 `json:"pending"`
	Inactive int64 `json:"inactive"`
}

type RecentSppgEntry struct {
	UserID          uint      `json:"user_id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	InstitutionName string    `json:"institution_name"`
	Capacity        int       `json:"capacity"`
	Address         string    `json:"address"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type RecentSchoolEntry struct {
	UserID       uint      `json:"user_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	NPSN         string    `json:"npsn"`
	Grade        string    `json:"grade"`
	TeacherCount int       `json:"teacher_count"`
	StudentCount int       `json:"student_count"`
	Address      string    `json:"address"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type RecentUmumEntry struct {
	UserID    uint      `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Lokasi    string    `json:"lokasi"`
	CreatedAt time.Time `json:"created_at"`
}

type DashboardData struct {
	RecentReports []FoodReport        `json:"recent_reports"`
	Sppg          []RecentSppgEntry   `json:"sppg"`
	SppgStatus    DashboardSppgStatus `json:"sppg_status"`
	Umum          []RecentUmumEntry   `json:"umum"`
	School        []RecentSchoolEntry `json:"school"`
	TotalSchool   int64               `json:"total_school"`
	TotalSppg     int64               `json:"total_sppg"`
	TotalStudent  int64               `json:"total_student"`
	TotalUmum     int64               `json:"total_umum"`
}

type DashboardResponse struct {
	Status string        `json:"status"`
	Data   DashboardData `json:"data"`
}
