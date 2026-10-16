package models

import "time"

// Jenis notifikasi. Dipakai klien untuk memilih ikon/warna, jadi nilainya
// stabil — jangan diubah tanpa menyesuaikan mobile.
const (
	NotifLaporanTerkirim    = "laporan_terkirim"    // konfirmasi ke pengirim laporan
	NotifLaporanMasuk       = "laporan_masuk"       // laporan baru untuk SPPG/sekolah
	NotifLaporanDitinjau    = "laporan_ditinjau"    // admin mengubah status laporan
	NotifPengirimanDatang   = "pengiriman_datang"   // SPPG mengirim makanan ke sekolah
	NotifPengirimanDiterima = "pengiriman_diterima" // sekolah konfirmasi terima
	NotifPeringatanAlergen  = "peringatan_alergen"  // hasil scan mengandung alergen
	NotifKemitraan          = "kemitraan"           // sekolah ditetapkan ke SPPG
	NotifPendaftaran        = "pendaftaran"         // status pendaftaran SPPG berubah
)

// Notification adalah satu baris di tab Notifikasi. Dibaca per user, jadi satu
// kejadian yang menyasar banyak orang menghasilkan banyak baris — itu disengaja
// supaya status baca/belum dibaca bisa berbeda per penerima.
type Notification struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null;index:idx_notif_user_created,priority:1" json:"user_id"`
	Type      string     `gorm:"type:varchar(40);not null" json:"type"`
	Title     string     `gorm:"type:varchar(150);not null" json:"title"`
	Body      string     `gorm:"type:text" json:"body"`
	RefType   string     `gorm:"type:varchar(40)" json:"ref_type"` // food_report, delivery_report, menu_report, sppg_registration
	RefID     uint       `json:"ref_id"`
	ReadAt    *time.Time `json:"read_at"` // null = belum dibaca
	CreatedAt time.Time  `gorm:"index:idx_notif_user_created,priority:2,sort:desc" json:"created_at"`
}

// IsRead dipakai klien supaya tidak perlu membandingkan null sendiri.
func (n Notification) IsRead() bool { return n.ReadAt != nil }
