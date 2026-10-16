package services

import (
	"log"

	"back-end/app/models"
	"back-end/config"
)

// Notify menyimpan satu notifikasi untuk satu user.
//
// Kegagalan tidak dikembalikan ke pemanggil: notifikasi adalah efek samping,
// dan laporan yang sudah tersimpan tidak boleh dianggap gagal hanya karena
// barisnya tidak masuk. Error tetap dicatat ke log supaya bisa ditelusuri.
func Notify(userID uint, notifType, title, body, refType string, refID uint) {
	if userID == 0 {
		return // tidak ada penerima (mis. sekolah belum punya mitra SPPG)
	}

	n := models.Notification{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Body:    body,
		RefType: refType,
		RefID:   refID,
	}

	if err := config.DB.Create(&n).Error; err != nil {
		log.Printf("[notif] gagal simpan notifikasi user=%d type=%s: %v", userID, notifType, err)
	}
}

// NotifyMany mengirim notifikasi yang sama ke beberapa user sekaligus.
// ID duplikat dan nol dibuang supaya satu user tidak menerima baris ganda.
func NotifyMany(userIDs []uint, notifType, title, body, refType string, refID uint) {
	if len(userIDs) == 0 {
		return
	}

	seen := make(map[uint]bool, len(userIDs))
	batch := make([]models.Notification, 0, len(userIDs))
	for _, id := range userIDs {
		if id == 0 || seen[id] {
			continue
		}
		seen[id] = true
		batch = append(batch, models.Notification{
			UserID:  id,
			Type:    notifType,
			Title:   title,
			Body:    body,
			RefType: refType,
			RefID:   refID,
		})
	}

	if len(batch) == 0 {
		return
	}

	if err := config.DB.Create(&batch).Error; err != nil {
		log.Printf("[notif] gagal simpan notifikasi massal type=%s n=%d: %v", notifType, len(batch), err)
	}
}

// NotifyAdmins mengirim notifikasi ke seluruh akun admin. Dipakai untuk
// kejadian yang perlu ditindaklanjuti BGN, mis. laporan masalah baru.
func NotifyAdmins(notifType, title, body, refType string, refID uint) {
	var adminIDs []uint
	if err := config.DB.Model(&models.User{}).
		Where("role_name = ?", "admin").
		Pluck("id", &adminIDs).Error; err != nil {
		log.Printf("[notif] gagal mengambil daftar admin: %v", err)
		return
	}
	NotifyMany(adminIDs, notifType, title, body, refType, refID)
}
